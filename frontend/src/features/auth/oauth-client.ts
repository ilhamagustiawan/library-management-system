import { createHash, randomBytes } from "node:crypto";

import { z } from "zod";

export type OAuthConfig = {
  issuer: string;
  serviceURL: string;
  clientId: string;
  clientSecret: string;
  redirectUri: string;
  scopes: string[];
};

export type OAuthFlow = {
  state: string;
  codeVerifier: string;
  createdAt: number;
};

export type OAuthTokens = {
  accessToken: string;
  refreshToken: string;
  tokenType: "Bearer";
  scope: string;
  expiresAt: number;
};

export type OAuthUser = {
  id: string;
  name: string;
  email: string;
};

type OAuthError = {
  kind: "unavailable" | "token-rejected" | "invalid-response";
};

export type TokenResult =
  | { status: "success"; tokens: OAuthTokens }
  | { status: "error"; error: OAuthError };

export type UserInfoResult =
  | { status: "success"; user: OAuthUser }
  | { status: "error"; error: OAuthError };

const tokenSchema = z.object({
  access_token: z.string().min(1),
  refresh_token: z.string().min(1),
  token_type: z.string().refine((value) => value.toLowerCase() === "bearer"),
  expires_in: z.number().int().positive(),
  scope: z.string(),
});

const userInfoSchema = z.object({
  data: z.object({
    id: z.string().min(1),
    name: z.string().min(1),
    email: z.email(),
  }),
});

function createFlow(now = Math.floor(Date.now() / 1_000)): OAuthFlow {
  return {
    state: randomBytes(32).toString("base64url"),
    codeVerifier: randomBytes(64).toString("base64url"),
    createdAt: now,
  };
}

function authorizeURL(config: OAuthConfig, flow: OAuthFlow) {
  const url = new URL("/oauth/authorize", config.issuer);
  url.searchParams.set("response_type", "code");
  url.searchParams.set("client_id", config.clientId);
  url.searchParams.set("redirect_uri", config.redirectUri);
  url.searchParams.set("scope", config.scopes.join(" "));
  url.searchParams.set("state", flow.state);
  url.searchParams.set(
    "code_challenge",
    createHash("sha256").update(flow.codeVerifier).digest("base64url"),
  );
  url.searchParams.set("code_challenge_method", "S256");
  return url;
}

async function requestTokens(
  config: OAuthConfig,
  body: URLSearchParams,
  fetcher: typeof fetch,
  now: number,
): Promise<TokenResult> {
  let response: Response;
  try {
    response = await fetcher(new URL("/oauth/token", config.serviceURL), {
      method: "POST",
      headers: {
        Accept: "application/json",
        Authorization: `Basic ${Buffer.from(`${config.clientId}:${config.clientSecret}`).toString("base64")}`,
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body,
      cache: "no-store",
    });
  } catch {
    return { status: "error", error: { kind: "unavailable" } };
  }
  if (!response.ok) {
    return { status: "error", error: { kind: "token-rejected" } };
  }

  try {
    const payload: unknown = await response.json();
    const parsed = tokenSchema.safeParse(payload);
    if (!parsed.success) {
      return { status: "error", error: { kind: "invalid-response" } };
    }
    return {
      status: "success",
      tokens: {
        accessToken: parsed.data.access_token,
        refreshToken: parsed.data.refresh_token,
        tokenType: "Bearer",
        scope: parsed.data.scope,
        expiresAt: now + parsed.data.expires_in,
      },
    };
  } catch {
    return { status: "error", error: { kind: "invalid-response" } };
  }
}

function exchangeCode(
  config: OAuthConfig,
  code: string,
  codeVerifier: string,
  fetcher: typeof fetch = fetch,
  now = Math.floor(Date.now() / 1_000),
) {
  return requestTokens(
    config,
    new URLSearchParams({
      grant_type: "authorization_code",
      code,
      redirect_uri: config.redirectUri,
      code_verifier: codeVerifier,
    }),
    fetcher,
    now,
  );
}

function refresh(
  config: OAuthConfig,
  refreshToken: string,
  fetcher: typeof fetch = fetch,
  now = Math.floor(Date.now() / 1_000),
) {
  return requestTokens(
    config,
    new URLSearchParams({ grant_type: "refresh_token", refresh_token: refreshToken }),
    fetcher,
    now,
  );
}

async function userInfo(
  config: OAuthConfig,
  accessToken: string,
  fetcher: typeof fetch = fetch,
): Promise<UserInfoResult> {
  let response: Response;
  try {
    response = await fetcher(new URL("/api/v1/oauth/userinfo", config.serviceURL), {
      headers: { Accept: "application/json", Authorization: `Bearer ${accessToken}` },
      cache: "no-store",
    });
  } catch {
    return { status: "error", error: { kind: "unavailable" } };
  }
  if (!response.ok) {
    return { status: "error", error: { kind: "token-rejected" } };
  }

  try {
    const payload: unknown = await response.json();
    const parsed = userInfoSchema.safeParse(payload);
    return parsed.success
      ? { status: "success", user: parsed.data.data }
      : { status: "error", error: { kind: "invalid-response" } };
  } catch {
    return { status: "error", error: { kind: "invalid-response" } };
  }
}

export const OAuthClient = { authorizeURL, createFlow, exchangeCode, refresh, userInfo } as const;
