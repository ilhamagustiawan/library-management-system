import { BookCheck, BookOpen, Clock3, ReceiptText } from "lucide-react";

import type { DashboardTab } from "./dashboard-tab";
import { LibraryDashboardTabs } from "./library-dashboard-tabs";
import type { LoadMemberLibraryResult } from "./member-library";

function formatMoney(amountMinor: number, currency: "IDR") {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amountMinor);
}

export function LibraryDashboard({
  initialTab,
  result,
}: {
  initialTab: DashboardTab;
  result: LoadMemberLibraryResult;
}) {
  if (result.status === "error") {
    const message = result.error.kind === "unauthorized"
      ? "Library access expired. Log out, then sign in again. Your account data remains unchanged."
      : "Your account remains signed in. Refresh this page; if the problem continues, check the library services.";
    return (
      <div role="alert" className="border border-destructive/40 bg-destructive/5 p-4 text-sm leading-6">
        <p className="font-semibold text-destructive">Library activity could not be loaded.</p>
        <p className="text-muted-foreground">{message}</p>
      </div>
    );
  }

  const { library } = result;
  const metrics = [
    { icon: BookOpen, label: "Checked out", value: String(library.summary.activeLoans), accent: "text-primary", surface: "bg-primary/5" },
    { icon: BookCheck, label: "Returned", value: String(library.summary.completedLoans), accent: "text-book-rust", surface: "bg-book-rust/5" },
    { icon: Clock3, label: "Late returns", value: String(library.summary.lateReturns), accent: "text-book-gold-dark", surface: "bg-book-gold/15" },
    { icon: ReceiptText, label: "Outstanding fines", value: formatMoney(library.summary.unpaidFineMinor, library.summary.fineCurrency), accent: "text-destructive", surface: "bg-destructive/5" },
  ] as const;

  return (
    <>
      <div className="grid overflow-hidden rounded-lg border border-border bg-card sm:grid-cols-2 xl:grid-cols-4">
        {metrics.map(({ icon: Icon, label, value, accent, surface }) => (
          <div key={label} className={`border-b border-border p-4 last:border-b-0 sm:nth-[2n+1]:border-r xl:border-b-0 xl:border-r xl:last:border-r-0 ${surface}`}>
            <Icon aria-hidden="true" className={`size-4 ${accent}`} strokeWidth={1.7} />
            <p className="mt-3 text-2xl font-semibold text-foreground">{value}</p>
            <p className="mt-1 text-xs font-semibold uppercase tracking-[0.08em] text-muted-foreground">{label}</p>
          </div>
        ))}
      </div>
      <LibraryDashboardTabs initialTab={initialTab} library={library} />
    </>
  );
}
