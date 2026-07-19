import { MemberHeader } from "@/components/member-header";
import type { DashboardTab } from "@/features/library/dashboard-tab";
import { LibraryDashboard } from "@/features/library/library-dashboard";
import type { DashboardNotice } from "@/features/library/dashboard-notice";
import { DashboardNoticeToast } from "@/features/library/dashboard-notice-toast";
import type { LoadMemberLibraryResult } from "@/features/library/member-library";

import type { OAuthUser } from "./oauth-client";

export function DashboardClient({
  library,
  initialTab,
  logoutEndpoint,
  notice,
  session,
}: {
  library: LoadMemberLibraryResult;
  initialTab: DashboardTab;
  logoutEndpoint: string;
  notice: DashboardNotice;
  session: OAuthUser;
}) {
  return (
    <div className="flex min-h-screen flex-col">
      {notice.kind === "book-borrowed" && <DashboardNoticeToast notice={notice} />}
      <MemberHeader logoutEndpoint={logoutEndpoint} />
      <main className="mx-auto w-full max-w-6xl flex-1 px-5 py-8 md:px-8 md:py-10">
        <div>
          <p className="text-xs font-bold uppercase tracking-[0.16em] text-book-rust">Member home</p>
          <h1 className="mt-2 text-3xl font-semibold leading-tight tracking-[-0.025em] sm:text-4xl">
            Welcome, {session.name}.
          </h1>
          <p className="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground sm:text-base">
            Current loans, return history, and fines—together in one clear record.
          </p>
        </div>
        <div className="mt-7">
          <LibraryDashboard initialTab={initialTab} result={library} />
        </div>
      </main>
      <footer className="border-t border-border bg-secondary px-5 py-5 text-center text-xs text-muted-foreground">
        Perpus Digital member portal
      </footer>
    </div>
  );
}
