import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { AuthRequest } from "@/features/auth/auth-request";

export async function POST(request: NextRequest) {
  const config = AuthConfig.load();
  if (!AuthRequest.hasSameOrigin(request)) {
    return NextResponse.json({ error: "invalid_origin" }, { status: 403 });
  }
  const response = NextResponse.redirect(new URL("/", request.url), 303);
  response.cookies.set(AuthCookies.sessionName, "", AuthCookies.options(config.secureCookies, 0));
  response.cookies.set(AuthCookies.flowName, "", AuthCookies.options(config.secureCookies, 0));
  response.headers.set("Cache-Control", "no-store");
  return response;
}
