import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

import { AuthCookies } from "@/features/auth/auth-cookies";
import { WebSession } from "@/features/auth/web-session";
import { POST } from "./route";

const appOrigin = "http://localhost:3000";
const sessionSecret = "abcdef0123456789abcdef0123456789";
const bookId = "0ec82798-8ff9-48c5-b68f-2b8c050647ac";

function sealedSession() {
  return WebSession.seal(
    {
      user: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
      accessToken: "access-token",
      refreshToken: "refresh-token",
      tokenType: "Bearer",
      scope: "books:read loans:borrow:self transactions:read:self",
      expiresAt: Math.floor(Date.now() / 1_000) + 900,
    },
    sessionSecret,
  );
}

function request(origin = appOrigin) {
  return new NextRequest(`${appOrigin}/api/member/loans`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      cookie: `${AuthCookies.sessionName}=${sealedSession()}`,
      origin,
    },
    body: JSON.stringify({ bookId }),
  });
}

describe("member loan route", () => {
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

  it("forwards only book ID with the server-held bearer token", async () => {
    const fetcher = vi.fn<typeof fetch>(async () => Response.json({ code: "LMS-200000", data: {} }, { status: 201 }));
    vi.stubGlobal("fetch", fetcher);

    const response = await POST(request());

    expect(response.status).toBe(201);
    const options = fetcher.mock.calls[0]?.[1];
    expect(options?.headers).toEqual(expect.objectContaining({ Authorization: "Bearer access-token" }));
    expect(options?.body).toBe(JSON.stringify({ bookId }));
  });

  it("maps the backend loan limit without exposing upstream details", async () => {
    vi.stubGlobal("fetch", vi.fn<typeof fetch>(async () => Response.json({ code: "LMS-409004", message: "internal" }, { status: 409 })));

    const response = await POST(request());

    expect(response.status).toBe(409);
    expect(await response.json()).toEqual({
      error: {
        kind: "loan-limit",
        message: "Loan limit reached. Return a book before borrowing another.",
      },
    });
  });

  it("rejects cross-origin requests before upstream access", async () => {
    const fetcher = vi.fn<typeof fetch>();
    vi.stubGlobal("fetch", fetcher);

    const response = await POST(request("https://attacker.test"));

    expect(response.status).toBe(403);
    expect(fetcher).not.toHaveBeenCalled();
  });
});
