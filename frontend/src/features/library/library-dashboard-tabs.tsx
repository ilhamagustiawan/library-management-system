"use client";

import { useEffect, useState, type KeyboardEvent } from "react";

import type { DashboardTab } from "./dashboard-tab";
import { DashboardTab as DashboardTabValue } from "./dashboard-tab";
import type { MemberLibrary } from "./member-library";
import { MemberBooksPanel, MemberHistoryPanel } from "./member-loan-panels";

function focusAdjacentTab(event: KeyboardEvent<HTMLButtonElement>) {
  if (!["ArrowLeft", "ArrowRight", "Home", "End"].includes(event.key)) return;
  event.preventDefault();
  const tabs = Array.from(
    event.currentTarget.parentElement?.querySelectorAll<HTMLButtonElement>("[role='tab']") ?? [],
  );
  const current = tabs.indexOf(event.currentTarget);
  const next = event.key === "Home"
    ? tabs.at(0)
    : event.key === "End"
      ? tabs.at(-1)
      : tabs.at((current + (event.key === "ArrowRight" ? 1 : -1) + tabs.length) % tabs.length);
  next?.focus();
}

export function LibraryDashboardTabs({
  initialTab,
  library,
}: {
  initialTab: DashboardTab;
  library: MemberLibrary;
}) {
  const [activeTab, setActiveTab] = useState(initialTab);

  useEffect(() => {
    function syncTab() {
      const params = new URLSearchParams(window.location.search);
      setActiveTab(DashboardTabValue.fromSearchParam(params.get("tab") ?? undefined));
    }
    window.addEventListener("popstate", syncTab);
    return () => window.removeEventListener("popstate", syncTab);
  }, []);

  function selectTab(tab: DashboardTab) {
    if (tab === activeTab) return;
    const params = new URLSearchParams(window.location.search);
    params.set("tab", tab);
    window.history.pushState(
      null,
      "",
      `${window.location.pathname}?${params.toString()}${window.location.hash}`,
    );
    setActiveTab(tab);
  }

  return (
    <div className="mt-5">
      <div
        aria-label="Dashboard sections"
        className="flex rounded-t-lg border border-border bg-card px-2"
        role="tablist"
      >
        <button
          aria-controls="books-panel"
          aria-selected={activeTab === "books"}
          className="-mb-px border-b-2 border-transparent px-4 py-3 text-sm font-semibold text-muted-foreground outline-none hover:text-primary focus-visible:ring-2 focus-visible:ring-ring data-[active=true]:border-primary data-[active=true]:text-primary"
          data-active={activeTab === "books"}
          id="books-tab"
          onClick={() => selectTab("books")}
          onKeyDown={focusAdjacentTab}
          role="tab"
          tabIndex={activeTab === "books" ? 0 : -1}
          type="button"
        >
          My Books
        </button>
        <button
          aria-controls="history-panel"
          aria-selected={activeTab === "history"}
          className="-mb-px border-b-2 border-transparent px-4 py-3 text-sm font-semibold text-muted-foreground outline-none hover:text-primary focus-visible:ring-2 focus-visible:ring-ring data-[active=true]:border-book-rust data-[active=true]:text-book-rust"
          data-active={activeTab === "history"}
          id="history-tab"
          onClick={() => selectTab("history")}
          onKeyDown={focusAdjacentTab}
          role="tab"
          tabIndex={activeTab === "history" ? 0 : -1}
          type="button"
        >
          History
        </button>
      </div>
      <div
        aria-labelledby="books-tab"
        hidden={activeTab !== "books"}
        id="books-panel"
        role="tabpanel"
        tabIndex={0}
      >
        <MemberBooksPanel loans={library.activeLoans} />
      </div>
      <div
        aria-labelledby="history-tab"
        hidden={activeTab !== "history"}
        id="history-panel"
        role="tabpanel"
        tabIndex={0}
      >
        <MemberHistoryPanel loans={library.history} />
      </div>
    </div>
  );
}
