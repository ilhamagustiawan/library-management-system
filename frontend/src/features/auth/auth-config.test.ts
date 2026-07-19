import { describe, expect, it } from "vitest";

import { AuthConfig } from "./auth-config";

describe("AuthConfig", () => {
  it("loads server-only OAuth and session settings", () => {
    const config = AuthConfig.load({
      AUTH_ISSUER: "http://localhost:8081",
      AUTH_CLIENT_ID: "nextjs",
      AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
      AUTH_REDIRECT_URI: "http://localhost:3000/api/auth/callback/library",
      AUTH_SCOPES: "library:read library:write",
      AUTH_SESSION_SECRET: "abcdef0123456789abcdef0123456789",
    });

    expect(config.oauth.scopes).toEqual(["library:read", "library:write"]);
    expect(config.secureCookies).toBe(false);
    expect(config.loginEndpoint).toBe("http://localhost:8081/api/v1/auth/login");
  });

  it("rejects weak session encryption secrets", () => {
    expect(() =>
      AuthConfig.load({
        AUTH_ISSUER: "http://localhost:8081",
        AUTH_CLIENT_ID: "nextjs",
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
        AUTH_CLIENT_ID: "nextjs",
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
        AUTH_ISSUER: "http://localhost:8081",
        AUTH_CLIENT_ID: "nextjs",
        AUTH_CLIENT_SECRET: "0123456789abcdef0123456789abcdef",
        AUTH_REDIRECT_URI: "http://localhost:3000/api/auth/callback/library",
        AUTH_SESSION_SECRET: "abcdef0123456789abcdef0123456789",
      }),
    ).toThrow("HTTPS");
  });
});
