import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { AuthRequest } from "@/features/auth/auth-request";
import { OAuthClient } from "@/features/auth/oauth-client";
import { OAuthFlowCookie } from "@/features/auth/oauth-flow-cookie";

export async function GET(request: NextRequest) {
  const config = AuthConfig.load();
  const callbackOrigin = new URL(config.oauth.redirectUri).origin;
  if (!AuthRequest.hasHost(request, callbackOrigin)) {
    return NextResponse.redirect(new URL("/api/auth/login", callbackOrigin));
  }
  const flow = OAuthClient.createFlow();
  const response = NextResponse.redirect(OAuthClient.authorizeURL(config.oauth, flow));
  response.cookies.set(
    AuthCookies.flowName,
    OAuthFlowCookie.seal(flow, config.sessionSecret),
    AuthCookies.options(config.secureCookies, OAuthFlowCookie.maxAgeSeconds),
  );
  response.headers.set("Cache-Control", "no-store");
  return response;
}
