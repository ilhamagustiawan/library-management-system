import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import { z } from "zod";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { AuthRequest } from "@/features/auth/auth-request";
import { WebSession } from "@/features/auth/web-session";

const inputSchema = z.object({ bookId: z.uuid() }).strict();
const upstreamSuccessSchema = z.object({ code: z.literal("LMS-200000") });
const upstreamErrorSchema = z.object({ code: z.string().min(1) });

type ErrorKind =
  | "invalid-input"
  | "session-expired"
  | "loan-limit"
  | "already-borrowed"
  | "unavailable"
  | "not-found"
  | "service-unavailable";

function error(status: number, kind: ErrorKind, message: string) {
  return NextResponse.json({ error: { kind, message } }, { status });
}

function upstreamError(status: number, code: string | undefined) {
  if (status === 401) return error(401, "session-expired", "Session expired. Log in again.");
  if (code === "LMS-409004") {
    return error(409, "loan-limit", "Loan limit reached. Return a book before borrowing another.");
  }
  if (code === "LMS-409005") {
    return error(409, "already-borrowed", "This book is already in your active loans.");
  }
  if (code === "LMS-409003") return error(409, "unavailable", "No copies currently available.");
  if (status === 404) return error(404, "not-found", "Book no longer exists in the catalog.");
  return error(503, "service-unavailable", "Book could not be borrowed. Try again.");
}

export async function POST(request: NextRequest) {
  if (!AuthRequest.hasSameOrigin(request)) {
    return error(403, "invalid-input", "Request origin rejected. Reload this page and try again.");
  }

  let body: unknown;
  try {
    body = await request.json();
  } catch {
    return error(422, "invalid-input", "Book selection is invalid. Return to Browse and select a book.");
  }
  const input = inputSchema.safeParse(body);
  if (!input.success) {
    return error(422, "invalid-input", "Book selection is invalid. Return to Browse and select a book.");
  }

  const config = AuthConfig.load();
  const sealedSession = request.cookies.get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) return error(401, "session-expired", "Session expired. Log in again.");
  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid" || WebSession.needsRefresh(openedSession.session)) {
    return error(401, "session-expired", "Session expired. Log in again.");
  }

  let response: Response;
  try {
    response = await fetch(new URL("/api/v1/transactions/loans", config.oauth.serviceURL), {
      method: "POST",
      headers: {
        Accept: "application/json",
        Authorization: `Bearer ${openedSession.session.accessToken}`,
        "Content-Type": "application/json",
      },
      cache: "no-store",
      body: JSON.stringify({ bookId: input.data.bookId }),
    });
  } catch {
    return error(503, "service-unavailable", "Loan service unavailable. Your account remains unchanged. Try again.");
  }

  let responseBody: unknown;
  try {
    responseBody = await response.json();
  } catch {
    return error(503, "service-unavailable", "Loan service returned an invalid response. Try again.");
  }
  if (response.status === 201) {
    const result = upstreamSuccessSchema.safeParse(responseBody);
    return result.success
      ? NextResponse.json({ status: "borrowed" }, { status: 201 })
      : error(503, "service-unavailable", "Loan service returned an invalid response. Try again.");
  }
  const result = upstreamErrorSchema.safeParse(responseBody);
  return upstreamError(response.status, result.success ? result.data.code : undefined);
}
