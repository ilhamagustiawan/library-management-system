import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

import { AuthCookies } from "@/features/auth/auth-cookies";
import { WebSession, type WebSession as WebSessionValue } from "@/features/auth/web-session";
import { POST as logout } from "./logout/route";
import { POST as refresh } from "./refresh/route";

const appOrigin = "http://localhost:3000";
const sessionSecret = "abcdef0123456789abcdef0123456789";

function session(overrides: Partial<WebSessionValue> = {}) {
  return WebSession.seal(
    {
      user: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
      accessToken: "access-token",
      refreshToken: "refresh-token",
      tokenType: "Bearer",
      scope: "books:read loans:borrow:self",
      expiresAt: Math.floor(Date.now() / 1_000) + 10,
      ...overrides,
    },
    sessionSecret,
  );
}

function request(
  path: string,
  origin: string,
  sealedSession?: string,
  serverOrigin = appOrigin,
) {
  const headers = new Headers({ origin });
  if (sealedSession !== undefined) {
    headers.set("cookie", `${AuthCookies.sessionName}=${sealedSession}`);
  }
  return new NextRequest(`${serverOrigin}${path}`, { method: "POST", headers });
}

function refreshedTokenResponse() {
  return Response.json({
    access_token: "new-access-token",
    refresh_token: "new-refresh-token",
    token_type: "Bearer",
    expires_in: 900,
    scope: "books:read loans:borrow:self",
  });
}

describe("session routes", () => {
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

  it("rejects cross-origin refresh", async () => {
    const response = await refresh(request("/api/auth/refresh", "https://attacker.test", session()));

    expect(response.status).toBe(403);
  });

  it("refreshes and rotates an expiring session cookie", async () => {
    vi.stubGlobal("fetch", vi.fn<typeof fetch>(async () => refreshedTokenResponse()));

    const response = await refresh(request("/api/auth/refresh", appOrigin, session()));

    expect(response.status).toBe(200);
    expect(await response.json()).toEqual({ status: "refreshed" });
    const sealed = response.cookies.get(AuthCookies.sessionName)?.value;
    expect(sealed).toBeDefined();
    if (sealed === undefined) return;
    const opened = WebSession.open(sealed, sessionSecret);
    expect(opened.status).toBe("valid");
    if (opened.status === "invalid") return;
    expect(opened.session.refreshToken).toBe("new-refresh-token");
  });

  it("shares one rotation between concurrent refresh requests", async () => {
    const fetcher = vi.fn<typeof fetch>(async () => refreshedTokenResponse());
    vi.stubGlobal("fetch", fetcher);
    const sealed = session({ refreshToken: "concurrent-refresh-token" });

    const [first, second] = await Promise.all([
      refresh(request("/api/auth/refresh", appOrigin, sealed)),
      refresh(request("/api/auth/refresh", appOrigin, sealed)),
    ]);

    expect(first.status).toBe(200);
    expect(second.status).toBe(200);
    expect(fetcher).toHaveBeenCalledOnce();
  });

  it("clears an expired session when rotation fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn<typeof fetch>(async () => Response.json({ error: "invalid_grant" }, { status: 400 })),
    );

    const response = await refresh(
      request("/api/auth/refresh", appOrigin, session({ refreshToken: "expired-refresh-token" })),
    );

    expect(response.status).toBe(401);
    expect(response.cookies.get(AuthCookies.sessionName)?.maxAge).toBe(0);
  });

  it("rejects cross-origin local logout", async () => {
    const response = await logout(request("/api/auth/logout", "https://attacker.test", session()));

    expect(response.status).toBe(403);
  });

  it("clears local cookies after same-origin logout", async () => {
    const response = await logout(request("/api/auth/logout", appOrigin, session()));

    expect(response.status).toBe(303);
    expect(response.cookies.get(AuthCookies.sessionName)?.maxAge).toBe(0);
    expect(response.cookies.get(AuthCookies.flowName)?.maxAge).toBe(0);
  });

  it("accepts configured origin when server uses an internal bind URL", async () => {
    const response = await logout(
      request("/api/auth/logout", appOrigin, session(), "http://0.0.0.0:3000"),
    );

    expect(response.status).toBe(303);
    expect(response.headers.get("location")).toBe(`${appOrigin}/`);
    expect(response.cookies.get(AuthCookies.sessionName)?.maxAge).toBe(0);
  });
});
