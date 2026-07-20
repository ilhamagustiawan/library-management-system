import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { AuthRequest } from "@/features/auth/auth-request";
import { OAuthClient } from "@/features/auth/oauth-client";
import { RefreshCoordinator } from "@/features/auth/refresh-coordinator";
import { WebSession } from "@/features/auth/web-session";

function expiredSession(secureCookies: boolean) {
  const response = NextResponse.json({ error: "session_expired" }, { status: 401 });
  response.cookies.set(AuthCookies.sessionName, "", AuthCookies.options(secureCookies, 0));
  response.headers.set("Cache-Control", "no-store");
  return response;
}

export async function POST(request: NextRequest) {
  const config = AuthConfig.load();
  if (!AuthRequest.hasOrigin(request, config.oauth.redirectUri)) {
    return NextResponse.json({ error: "invalid_origin" }, { status: 403 });
  }
  const sealedSession = request.cookies.get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) {
    return expiredSession(config.secureCookies);
  }
  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid") {
    return expiredSession(config.secureCookies);
  }
  if (!WebSession.needsRefresh(openedSession.session)) {
    const response = NextResponse.json({ status: "active" });
    response.headers.set("Cache-Control", "no-store");
    return response;
  }

  const refreshed = await RefreshCoordinator.run(openedSession.session.refreshToken, () =>
    OAuthClient.refresh(config.oauth, openedSession.session.refreshToken),
  );
  if (refreshed.status === "error") {
    return expiredSession(config.secureCookies);
  }

  const response = NextResponse.json({ status: "refreshed" });
  response.cookies.set(
    AuthCookies.sessionName,
    WebSession.seal(
      {
        user: openedSession.session.user,
        accessToken: refreshed.tokens.accessToken,
        refreshToken: refreshed.tokens.refreshToken,
        tokenType: refreshed.tokens.tokenType,
        scope: refreshed.tokens.scope,
        expiresAt: refreshed.tokens.expiresAt,
      },
      config.sessionSecret,
    ),
    AuthCookies.options(config.secureCookies, AuthCookies.sessionMaxAgeSeconds),
  );
  response.headers.set("Cache-Control", "no-store");
  return response;
}
