import type { Metadata } from "next";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { DashboardClient } from "@/features/auth/dashboard-client";
import { SessionRefresh } from "@/features/auth/session-refresh";
import { WebSession } from "@/features/auth/web-session";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Member home",
};

export default async function DashboardPage() {
  const config = AuthConfig.load();
  const sealedSession = (await cookies()).get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) redirect("/login");

  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid") redirect("/login");
  if (WebSession.needsRefresh(openedSession.session)) {
    return <SessionRefresh />;
  }

  return <DashboardClient logoutEndpoint={config.logoutEndpoint} session={openedSession.session.user} />;
}
