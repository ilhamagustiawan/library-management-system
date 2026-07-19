import { z } from "zod";

import type { OAuthConfig } from "./oauth-client";

const envSchema = z.object({
  NODE_ENV: z.enum(["development", "test", "production"]).default("development"),
  AUTH_ISSUER: z.url(),
  USER_SERVICE_URL: z.url().default("http://localhost:8000"),
  AUTH_CLIENT_ID: z.string().min(1),
  AUTH_CLIENT_SECRET: z.string().min(32).max(72),
  AUTH_REDIRECT_URI: z.url(),
  AUTH_SCOPES: z
    .string()
    .default("loans:borrow:self loans:return:self transactions:read:self books:read"),
  AUTH_SESSION_SECRET: z.string().min(32),
});

export type AuthConfig = {
  oauth: OAuthConfig;
  sessionSecret: string;
  secureCookies: boolean;
  loginEndpoint: string;
  logoutEndpoint: string;
  registerEndpoint: string;
};

const loopbackHosts = new Set(["localhost", "127.0.0.1", "[::1]"]);

function validateTransport(url: URL, field: string, environment: string) {
  if (url.protocol === "https:") return;
  if (environment !== "production" && url.protocol === "http:" && loopbackHosts.has(url.hostname)) {
    return;
  }
  throw new Error(`Invalid auth configuration: ${field} must use HTTPS except loopback development`);
}

function load(environment: Record<string, string | undefined> = process.env): AuthConfig {
  const result = envSchema.safeParse(environment);
  if (!result.success) {
    const fields = [...new Set(result.error.issues.map((issue) => issue.path.join(".")))].join(", ");
    throw new Error(`Invalid auth configuration: ${fields}`);
  }

  const issuer = new URL(result.data.AUTH_ISSUER);
  const userService = new URL(result.data.USER_SERVICE_URL);
  const redirect = new URL(result.data.AUTH_REDIRECT_URI);
  if (issuer.pathname !== "/" || issuer.search !== "" || issuer.hash !== "") {
    throw new Error("Invalid auth configuration: AUTH_ISSUER must be an origin");
  }
  if (userService.pathname !== "/" || userService.search !== "" || userService.hash !== "") {
    throw new Error("Invalid auth configuration: USER_SERVICE_URL must be an origin");
  }
  if (redirect.hash !== "") {
    throw new Error("Invalid auth configuration: AUTH_REDIRECT_URI must not contain a fragment");
  }
  validateTransport(issuer, "AUTH_ISSUER", result.data.NODE_ENV);
  validateTransport(userService, "USER_SERVICE_URL", result.data.NODE_ENV);
  validateTransport(redirect, "AUTH_REDIRECT_URI", result.data.NODE_ENV);
  const scopes = result.data.AUTH_SCOPES.split(/\s+/).filter(Boolean);
  if (scopes.length === 0) {
    throw new Error("Invalid auth configuration: AUTH_SCOPES");
  }

  return {
    oauth: {
      issuer: issuer.origin,
      clientId: result.data.AUTH_CLIENT_ID,
      clientSecret: result.data.AUTH_CLIENT_SECRET,
      redirectUri: redirect.toString(),
      scopes,
    },
    sessionSecret: result.data.AUTH_SESSION_SECRET,
    secureCookies: redirect.protocol === "https:",
    loginEndpoint: new URL("/api/v1/auth/login", issuer).toString(),
    logoutEndpoint: new URL("/api/v1/auth/logout", issuer).toString(),
    registerEndpoint: new URL("/api/v1/users", userService).toString(),
  };
}

export const AuthConfig = { load } as const;
