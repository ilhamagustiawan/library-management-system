import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";
import { NextRequest } from "next/server";

import { AuthCookies } from "@/features/auth/auth-cookies";
import type { OAuthFlow } from "@/features/auth/oauth-client";
import { OAuthFlowCookie } from "@/features/auth/oauth-flow-cookie";
import { GET as callback } from "./callback/library/route";
import { GET as startLogin } from "./login/route";

const appOrigin = "http://localhost:3000";
const authOrigin = "http://localhost:8000";
const sessionSecret = "abcdef0123456789abcdef0123456789";

function callbackRequest(query: string, flow?: OAuthFlow, serverOrigin = appOrigin) {
  const headers = new Headers({ host: new URL(appOrigin).host });
  if (flow !== undefined) {
    headers.set(
      "cookie",
      `${AuthCookies.flowName}=${OAuthFlowCookie.seal(flow, sessionSecret)}`,
    );
  }
  return new NextRequest(`${serverOrigin}/api/auth/callback/library?${query}`, { headers });
}

function tokenResponse() {
  return Response.json({
    access_token: "access-token",
    refresh_token: "refresh-token",
    token_type: "Bearer",
    expires_in: 900,
    scope: "books:read loans:borrow:self",
  });
}

describe("OAuth login routes", () => {
  const flow: OAuthFlow = {
    state: "expected-state",
    codeVerifier: "code-verifier",
    createdAt: Math.floor(Date.now() / 1_000),
  };

  beforeEach(() => {
    vi.stubEnv("AUTH_ISSUER", authOrigin);
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", `${appOrigin}/api/auth/callback/library`);
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    vi.unstubAllGlobals();
  });

  it("starts authorization with PKCE and a sealed flow cookie", async () => {
    const response = await startLogin(new NextRequest(`${appOrigin}/api/auth/login`));

    const location = response.headers.get("location");
    expect(location).not.toBeNull();
    if (location === null) return;
    const target = new URL(location);
    expect(target.origin).toBe(authOrigin);
    expect(target.searchParams.get("code_challenge_method")).toBe("S256");
    expect(response.cookies.get(AuthCookies.flowName)?.httpOnly).toBe(true);
  });

  it("moves login to the configured callback origin before setting flow state", async () => {
    const response = await startLogin(
      new NextRequest("http://0.0.0.0:3000/api/auth/login"),
    );

    expect(response.headers.get("location")).toBe(`${appOrigin}/api/auth/login`);
    expect(response.cookies.get(AuthCookies.flowName)).toBeUndefined();
  });

  it("starts login when public host matches despite internal bind URL", async () => {
    const response = await startLogin(
      new NextRequest("http://0.0.0.0:3000/api/auth/login", {
        headers: { host: "localhost:3000" },
      }),
    );

    expect(response.headers.get("location")).toContain(`${authOrigin}/oauth/authorize`);
    expect(response.cookies.get(AuthCookies.flowName)?.httpOnly).toBe(true);
  });

  it("rejects callbacks without a flow cookie", async () => {
    const response = await callback(callbackRequest("code=code-123&state=expected-state"));

    expect(response.headers.get("location")).toContain("error=invalid_callback");
    expect(response.cookies.get(AuthCookies.flowName)?.maxAge).toBe(0);
  });

  it("rejects mismatched callback state before token exchange", async () => {
    const fetcher = vi.fn<typeof fetch>();
    vi.stubGlobal("fetch", fetcher);

    const response = await callback(callbackRequest("code=code-123&state=wrong", flow));

    expect(response.headers.get("location")).toContain("error=invalid_callback");
    expect(fetcher).not.toHaveBeenCalled();
  });

  it("reports token exchange failure without creating a session", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn<typeof fetch>(async () => Response.json({ error: "invalid_grant" }, { status: 400 })),
    );

    const response = await callback(
      callbackRequest(
        "code=code-123&state=expected-state",
        flow,
        "http://0.0.0.0:3000",
      ),
    );

    expect(response.headers.get("location")).toBe(
      `${appOrigin}/login?error=token_exchange_failed`,
    );
    expect(response.cookies.get(AuthCookies.sessionName)).toBeUndefined();
  });

  it("reports user-info failure without creating a session", async () => {
    const fetcher = vi.fn<typeof fetch>();
    fetcher.mockResolvedValueOnce(tokenResponse());
    fetcher.mockResolvedValueOnce(Response.json({ error: "unavailable" }, { status: 503 }));
    vi.stubGlobal("fetch", fetcher);

    const response = await callback(
      callbackRequest(
        "code=code-123&state=expected-state",
        flow,
        "http://0.0.0.0:3000",
      ),
    );

    expect(response.headers.get("location")).toContain("error=user_info_failed");
    expect(response.cookies.get(AuthCookies.sessionName)).toBeUndefined();
  });

  it("creates an encrypted HttpOnly session after successful callback", async () => {
    const fetcher = vi.fn<typeof fetch>();
    fetcher.mockResolvedValueOnce(tokenResponse());
    fetcher.mockResolvedValueOnce(
      Response.json({
        code: "LMS-200000",
        data: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
      }),
    );
    vi.stubGlobal("fetch", fetcher);

    const response = await callback(
      callbackRequest(
        "code=code-123&state=expected-state",
        flow,
        "http://0.0.0.0:3000",
      ),
    );

    expect(response.headers.get("location")).toBe(`${appOrigin}/dashboard`);
    const sessionCookie = response.cookies.get(AuthCookies.sessionName);
    expect(sessionCookie?.httpOnly).toBe(true);
    expect(sessionCookie?.value).not.toContain("access-token");
    expect(response.cookies.get(AuthCookies.flowName)?.maxAge).toBe(0);
  });
});
