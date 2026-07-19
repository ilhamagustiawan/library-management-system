import { NextResponse } from "next/server";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { OAuthClient } from "@/features/auth/oauth-client";
import { OAuthFlowCookie } from "@/features/auth/oauth-flow-cookie";

export async function GET() {
  const config = AuthConfig.load();
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
