import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { z } from "zod";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { AuthRequest } from "@/features/auth/auth-request";
import { WebSession } from "@/features/auth/web-session";

const paramsSchema = z.object({ loanId: z.uuid() });
const returnInputSchema = z.object({ acceptedFineAmountMinor: z.number().int().nonnegative() }).strict();
const quoteFineSchema = z.object({
  overdueDays: z.number().int().positive(),
  dailyRateMinor: z.number().int().positive(),
  totalAmountMinor: z.number().int().positive(),
  currency: z.literal("IDR"),
});
const quoteSuccessSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    loanId: z.uuid(),
    bookId: z.uuid(),
    dueAt: z.iso.datetime(),
    quotedAt: z.iso.datetime(),
    fine: quoteFineSchema.nullable(),
  }),
});
const returnSuccessSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    status: z.literal("returned"),
    stockSyncStatus: z.enum(["pending", "confirmed"]),
    fine: z
      .object({
        overdueDays: z.number().int().positive(),
        totalAmountMinor: z.number().int().positive(),
        currency: z.literal("IDR"),
      })
      .optional(),
  }),
});
const upstreamErrorSchema = z.object({ code: z.string().min(1) });

type ErrorKind =
  | "invalid-input"
  | "session-expired"
  | "forbidden"
  | "not-found"
  | "fine-quote-changed"
  | "service-unavailable";

type Context = { params: Promise<{ loanId: string }> };

function error(status: number, kind: ErrorKind, message: string) {
  return NextResponse.json({ error: { kind, message } }, { status });
}

function mapUpstreamError(status: number, code: string | undefined) {
  if (status === 401) return error(401, "session-expired", "Session expired. Log in again.");
  if (status === 403) {
    return error(403, "forbidden", "This loan cannot be returned from your account.");
  }
  if (status === 404) {
    return error(404, "not-found", "Active loan no longer exists. Refresh your library.");
  }
  if (status === 409 && code === "LMS-409006") {
    return error(
      409,
      "fine-quote-changed",
      "Fine changed before return. Review the updated amount and confirm again.",
    );
  }
  return error(
    503,
    "service-unavailable",
    "Loan service unavailable. Your loan remains unchanged. Try again.",
  );
}

function accessToken(request: NextRequest) {
  const config = AuthConfig.load();
  const sealedSession = request.cookies.get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) {
    return { status: "error", response: error(401, "session-expired", "Session expired. Log in again.") } as const;
  }
  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid" || WebSession.needsRefresh(openedSession.session)) {
    return { status: "error", response: error(401, "session-expired", "Session expired. Log in again.") } as const;
  }
  return {
    status: "success",
    issuer: config.oauth.serviceURL,
    accessToken: openedSession.session.accessToken,
  } as const;
}

async function responseBody(response: Response): Promise<unknown | undefined> {
  try {
    const body: unknown = await response.json();
    return body;
  } catch {
    return undefined;
  }
}

async function loanID(context: Context) {
  return paramsSchema.safeParse(await context.params);
}

export async function GET(request: NextRequest, context: Context) {
  const params = await loanID(context);
  if (!params.success) return error(422, "invalid-input", "Loan selection is invalid.");
  const session = accessToken(request);
  if (session.status === "error") return session.response;

  let response: Response;
  try {
    response = await fetch(
      new URL(`/api/v1/transactions/loans/${encodeURIComponent(params.data.loanId)}/return`, session.issuer),
      {
        headers: { Accept: "application/json", Authorization: `Bearer ${session.accessToken}` },
        cache: "no-store",
      },
    );
  } catch {
    return mapUpstreamError(503, undefined);
  }
  const body = await responseBody(response);
  if (!response.ok) {
    const upstreamError = upstreamErrorSchema.safeParse(body);
    return mapUpstreamError(response.status, upstreamError.success ? upstreamError.data.code : undefined);
  }
  const result = quoteSuccessSchema.safeParse(body);
  return result.success
    ? NextResponse.json({ status: "ready", quote: result.data.data })
    : error(503, "service-unavailable", "Loan service returned an invalid quote. Try again.");
}

export async function POST(request: NextRequest, context: Context) {
  const config = AuthConfig.load();
  if (!AuthRequest.hasOrigin(request, config.oauth.redirectUri)) {
    return error(403, "invalid-input", "Request origin rejected. Reload this page and try again.");
  }
  const params = await loanID(context);
  if (!params.success) return error(422, "invalid-input", "Loan selection is invalid.");
  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return error(422, "invalid-input", "Fine confirmation is invalid. Review the quote again.");
  }
  const input = returnInputSchema.safeParse(body);
  if (!input.success) {
    return error(422, "invalid-input", "Fine confirmation is invalid. Review the quote again.");
  }
  const session = accessToken(request);
  if (session.status === "error") return session.response;

  let response: Response;
  try {
    response = await fetch(
      new URL(`/api/v1/transactions/loans/${encodeURIComponent(params.data.loanId)}/return`, session.issuer),
      {
        method: "POST",
        headers: {
          Accept: "application/json",
          Authorization: `Bearer ${session.accessToken}`,
          "Content-Type": "application/json",
        },
        cache: "no-store",
        body: JSON.stringify(input.data),
      },
    );
  } catch {
    return mapUpstreamError(503, undefined);
  }
  const upstreamBody = await responseBody(response);
  if (response.status !== 200 && response.status !== 202) {
    const upstreamError = upstreamErrorSchema.safeParse(upstreamBody);
    return mapUpstreamError(response.status, upstreamError.success ? upstreamError.data.code : undefined);
  }
  const result = returnSuccessSchema.safeParse(upstreamBody);
  if (!result.success) {
    return error(503, "service-unavailable", "Loan service returned an invalid result. Refresh your library.");
  }
  const fine = result.data.data.fine;
  return NextResponse.json(
    {
      status: "returned",
      stockUpdate: response.status === 202 ? "pending" : "confirmed",
      fine: fine === undefined
        ? null
        : {
            overdueDays: fine.overdueDays,
            totalAmountMinor: fine.totalAmountMinor,
            currency: fine.currency,
          },
    },
    { status: response.status },
  );
}
