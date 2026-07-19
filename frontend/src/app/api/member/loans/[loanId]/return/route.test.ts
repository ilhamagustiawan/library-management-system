import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

import { AuthCookies } from "@/features/auth/auth-cookies";
import { WebSession } from "@/features/auth/web-session";

import { GET, POST } from "./route";

const appOrigin = "http://localhost:3000";
const sessionSecret = "abcdef0123456789abcdef0123456789";
const loanId = "52a88672-a4c2-4876-be5a-65863aeb35e4";

function sealedSession() {
  return WebSession.seal(
    {
      user: { id: "member-1", name: "Maya Chen", email: "maya@perpus-digital.test" },
      accessToken: "access-token",
      refreshToken: "refresh-token",
      tokenType: "Bearer",
      scope: "loans:return:self transactions:read:self",
      expiresAt: Math.floor(Date.now() / 1_000) + 900,
    },
    sessionSecret,
  );
}

function request(method: "GET" | "POST", body?: string, origin = appOrigin) {
  return new NextRequest(`${appOrigin}/api/member/loans/${loanId}/return`, {
    method,
    headers: {
      cookie: `${AuthCookies.sessionName}=${sealedSession()}`,
      ...(method === "POST" ? { "Content-Type": "application/json", origin } : {}),
    },
    body,
  });
}

const context = { params: Promise.resolve({ loanId }) };

describe("member return route", () => {
  beforeEach(() => {
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8000");
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", `${appOrigin}/api/auth/callback/library`);
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
  });

  it("loads an authoritative fine quote with the server-held token", async () => {
    const fetcher = vi.fn<typeof fetch>(async () =>
      Response.json({
        code: "LMS-200000",
        data: {
          loanId,
          bookId: "7b36fe43-f31d-4861-884f-42ed7386b1e9",
          dueAt: "2026-07-26T10:00:00Z",
          quotedAt: "2026-07-28T10:00:00Z",
          fine: {
            overdueDays: 2,
            dailyRateMinor: 5000,
            totalAmountMinor: 10000,
            currency: "IDR",
          },
        },
      }),
    );
    vi.stubGlobal("fetch", fetcher);

    const response = await GET(request("GET"), context);

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({
      status: "ready",
      quote: expect.objectContaining({ loanId, fine: expect.objectContaining({ totalAmountMinor: 10000 }) }),
    });
    expect(fetcher.mock.calls[0]?.[1]?.headers).toEqual(
      expect.objectContaining({ Authorization: "Bearer access-token" }),
    );
  });

  it("returns with the fine amount accepted by the member", async () => {
    const fetcher = vi.fn<typeof fetch>(async () =>
      Response.json(
        {
          code: "LMS-200000",
          data: {
            status: "returned",
            stockSyncStatus: "pending",
            fine: {
              overdueDays: 2,
              dailyRateMinor: 5000,
              totalAmountMinor: 10000,
              currency: "IDR",
              status: "unpaid",
            },
          },
        },
        { status: 202 },
      ),
    );
    vi.stubGlobal("fetch", fetcher);

    const response = await POST(
      request("POST", JSON.stringify({ acceptedFineAmountMinor: 10000 })),
      context,
    );

    expect(response.status).toBe(202);
    expect(await response.json()).toEqual({
      status: "returned",
      stockUpdate: "pending",
      fine: { overdueDays: 2, totalAmountMinor: 10000, currency: "IDR" },
    });
    expect(fetcher.mock.calls[0]?.[1]?.body).toBe(JSON.stringify({ acceptedFineAmountMinor: 10000 }));
  });

  it("rejects cross-origin returns before upstream access", async () => {
    const fetcher = vi.fn<typeof fetch>();
    vi.stubGlobal("fetch", fetcher);

    const response = await POST(
      request("POST", JSON.stringify({ acceptedFineAmountMinor: 0 }), "https://attacker.test"),
      context,
    );

    expect(response.status).toBe(403);
    expect(fetcher).not.toHaveBeenCalled();
  });

  it("maps a changed fine quote to a typed conflict", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn<typeof fetch>(async () =>
        Response.json({ code: "LMS-409006", message: "internal" }, { status: 409 }),
      ),
    );

    const response = await POST(
      request("POST", JSON.stringify({ acceptedFineAmountMinor: 0 })),
      context,
    );

    expect(response.status).toBe(409);
    expect(await response.json()).toEqual({
      error: {
        kind: "fine-quote-changed",
        message: "Fine changed before return. Review the updated amount and confirm again.",
      },
    });
  });
});
