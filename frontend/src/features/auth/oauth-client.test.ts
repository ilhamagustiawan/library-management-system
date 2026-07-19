import { createHash } from "node:crypto";

import { describe, expect, it } from "vitest";

import { OAuthClient, type OAuthConfig } from "./oauth-client";

const config: OAuthConfig = {
  issuer: "http://localhost:8081",
  clientId: "nextjs",
  clientSecret: "0123456789abcdef0123456789abcdef",
  redirectUri: "http://localhost:3000/api/auth/callback/library",
  scopes: ["library:read", "library:write"],
};

describe("OAuthClient", () => {
  it("creates an authorization request with state and S256 PKCE", () => {
    const flow = OAuthClient.createFlow();
    const authorizeURL = OAuthClient.authorizeURL(config, flow);

    expect(authorizeURL.searchParams.get("state")).toBe(flow.state);
    expect(authorizeURL.searchParams.get("code_challenge_method")).toBe("S256");
    expect(authorizeURL.searchParams.get("code_challenge")).toBe(
      createHash("sha256").update(flow.codeVerifier).digest("base64url"),
    );
  });

  it("rotates a refresh token and returns normalized expiry", async () => {
    const fetcher: typeof fetch = async (input, init) => {
      expect(String(input)).toBe("http://localhost:8081/oauth/token");
      expect(init?.headers).toEqual(
        expect.objectContaining({ Authorization: expect.stringMatching(/^Basic /) }),
      );
      expect(String(init?.body)).toContain("grant_type=refresh_token");
      expect(String(init?.body)).toContain("refresh_token=old-refresh-token");
      return Response.json({
        access_token: "new-access-token",
        refresh_token: "new-refresh-token",
        token_type: "Bearer",
        expires_in: 900,
        scope: "library:read library:write",
      });
    };

    const result = await OAuthClient.refresh(config, "old-refresh-token", fetcher, 1_000);

    expect(result).toEqual({
      status: "success",
      tokens: {
        accessToken: "new-access-token",
        refreshToken: "new-refresh-token",
        tokenType: "Bearer",
        scope: "library:read library:write",
        expiresAt: 1_900,
      },
    });
  });

  it("returns a structured error for rejected refresh", async () => {
    const fetcher: typeof fetch = async () =>
      Response.json({ error: "invalid_grant" }, { status: 400 });

    await expect(OAuthClient.refresh(config, "expired", fetcher, 1_000)).resolves.toEqual({
      status: "error",
      error: { kind: "token-rejected" },
    });
  });
});
