import { describe, expect, it } from "vitest";

import { AuthConfig } from "./auth-config";

describe("AuthConfig", () => {
  it("loads server-only OAuth and session settings", () => {
    const config = AuthConfig.load({
      AUTH_ISSUER: "http://localhost:8000",
      USER_SERVICE_URL: "http://localhost:8000",
      AUTH_CLIENT_ID: "member-nextjs-web",
      AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
      AUTH_REDIRECT_URI: "http://localhost:3000/api/auth/callback/library",
      AUTH_SCOPES: "loans:borrow:self loans:return:self transactions:read:self books:read",
      AUTH_SESSION_SECRET: "abcdef0123456789abcdef0123456789",
    });

    expect(config.oauth.scopes).toEqual([
      "loans:borrow:self",
      "loans:return:self",
      "transactions:read:self",
      "books:read",
    ]);
    expect(config.oauth.clientId).toBe("member-nextjs-web");
    expect(config.secureCookies).toBe(false);
    expect(config.loginEndpoint).toBe("http://localhost:8000/api/v1/auth/login");
    expect(config.registerEndpoint).toBe("http://localhost:8000/api/v1/users");
  });

  it("rejects weak session encryption secrets", () => {
    expect(() =>
      AuthConfig.load({
        AUTH_ISSUER: "http://localhost:8000",
        AUTH_CLIENT_ID: "member-nextjs-web",
        AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
        AUTH_REDIRECT_URI: "http://localhost:3000/api/auth/callback/library",
        AUTH_SESSION_SECRET: "too-short",
      }),
    ).toThrow("AUTH_SESSION_SECRET");
  });

  it("rejects insecure non-loopback auth URLs", () => {
    expect(() =>
      AuthConfig.load({
        AUTH_ISSUER: "http://auth.example",
        AUTH_CLIENT_ID: "member-nextjs-web",
        AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
        AUTH_REDIRECT_URI: "http://app.example/api/auth/callback/library",
        AUTH_SESSION_SECRET: "abcdef0123456789abcdef0123456789",
      }),
    ).toThrow("HTTPS");
  });

  it("rejects loopback HTTP in production", () => {
    expect(() =>
      AuthConfig.load({
        NODE_ENV: "production",
        AUTH_ISSUER: "http://localhost:8000",
        AUTH_CLIENT_ID: "member-nextjs-web",
        AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
        AUTH_REDIRECT_URI: "http://localhost:3000/api/auth/callback/library",
        AUTH_SESSION_SECRET: "abcdef0123456789abcdef0123456789",
      }),
    ).toThrow("HTTPS");
  });
});
