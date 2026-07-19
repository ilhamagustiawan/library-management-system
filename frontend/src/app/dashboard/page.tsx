import type { Metadata } from "next";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { DashboardClient } from "@/features/auth/dashboard-client";
import { SessionRefresh } from "@/features/auth/session-refresh";
import { WebSession } from "@/features/auth/web-session";
import { MemberLibrary } from "@/features/library/member-library";
import { DashboardNotice } from "@/features/library/dashboard-notice";
import { DashboardTab } from "@/features/library/dashboard-tab";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Member home",
};

type DashboardPageProps = {
  searchParams: Promise<{
    borrowed?: string | string[];
    tab?: string | string[];
  }>;
};

export default async function DashboardPage({ searchParams }: DashboardPageProps) {
  const config = AuthConfig.load();
  const sealedSession = (await cookies()).get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) redirect("/login");

  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid") redirect("/login");
  if (WebSession.needsRefresh(openedSession.session)) {
    return <SessionRefresh />;
  }

  const library = await MemberLibrary.load({
    issuer: config.oauth.issuer,
    accessToken: openedSession.session.accessToken,
  });
  const params = await searchParams;
  const notice = DashboardNotice.fromSearchParams(params);

  return (
    <DashboardClient
      library={library}
      initialTab={DashboardTab.fromSearchParam(params.tab)}
      logoutEndpoint={config.logoutEndpoint}
      notice={notice}
      session={openedSession.session.user}
    />
  );
}
