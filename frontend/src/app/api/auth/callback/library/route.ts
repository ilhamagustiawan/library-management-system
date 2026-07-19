import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { OAuthCallback } from "@/features/auth/oauth-callback";
import { OAuthClient } from "@/features/auth/oauth-client";
import { OAuthFlowCookie } from "@/features/auth/oauth-flow-cookie";
import { WebSession } from "@/features/auth/web-session";

function loginError(request: NextRequest, secureCookies: boolean, error: string) {
  const url = new URL("/login", request.url);
  url.searchParams.set("error", error);
  const response = NextResponse.redirect(url);
  response.cookies.set(AuthCookies.flowName, "", AuthCookies.options(secureCookies, 0));
  response.headers.set("Cache-Control", "no-store");
  return response;
}

export async function GET(request: NextRequest) {
  const config = AuthConfig.load();
  const sealedFlow = request.cookies.get(AuthCookies.flowName)?.value;
  if (sealedFlow === undefined) {
    return loginError(request, config.secureCookies, "invalid_callback");
  }
  const openedFlow = OAuthFlowCookie.open(sealedFlow, config.sessionSecret);
  if (openedFlow.status === "invalid") {
    return loginError(request, config.secureCookies, "invalid_callback");
  }

  const callback = OAuthCallback.validate(
    request.nextUrl.searchParams.get("code"),
    request.nextUrl.searchParams.get("state"),
    request.nextUrl.searchParams.get("error"),
    openedFlow.flow,
  );
  if (callback.status === "error") {
    return loginError(request, config.secureCookies, callback.error);
  }

  const tokenResult = await OAuthClient.exchangeCode(
    config.oauth,
    callback.code,
    openedFlow.flow.codeVerifier,
  );
  if (tokenResult.status === "error") {
    return loginError(request, config.secureCookies, "token_exchange_failed");
  }
  const userResult = await OAuthClient.userInfo(config.oauth, tokenResult.tokens.accessToken);
  if (userResult.status === "error") {
    return loginError(request, config.secureCookies, "user_info_failed");
  }

  const response = NextResponse.redirect(new URL("/dashboard", request.url));
  response.cookies.set(
    AuthCookies.sessionName,
    WebSession.seal(
      {
        user: userResult.user,
        accessToken: tokenResult.tokens.accessToken,
        refreshToken: tokenResult.tokens.refreshToken,
        tokenType: tokenResult.tokens.tokenType,
        scope: tokenResult.tokens.scope,
        expiresAt: tokenResult.tokens.expiresAt,
      },
      config.sessionSecret,
    ),
    AuthCookies.options(config.secureCookies, AuthCookies.sessionMaxAgeSeconds),
  );
  response.cookies.set(AuthCookies.flowName, "", AuthCookies.options(config.secureCookies, 0));
  response.headers.set("Cache-Control", "no-store");
  return response;
}
