import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { DashboardClient } from "./dashboard-client";
import type { MemberLibrary } from "@/features/library/member-library";

vi.mock("next/navigation", () => ({
  usePathname: () => "/dashboard",
  useRouter: () => ({ replace: vi.fn() }),
}));

const session = {
  id: "user-123",
  name: "Maya Chen",
  email: "maya@perpus-digital.test",
};

const library: MemberLibrary = {
  activeLoans: [
    {
      loanId: "loan-active",
      book: {
        status: "available",
        id: "book-active",
        title: "Atomic Habits",
        author: "James Clear",
        coverUrl: null,
        publicationYear: 2016,
      },
      borrowedAt: "2026-07-17T10:00:00Z",
      state: { status: "active" },
    },
  ],
  history: [
    {
      loanId: "loan-late",
      book: {
        status: "available",
        id: "book-late",
        title: "The Alchemist",
        author: "Paulo Coelho",
        coverUrl: null,
        publicationYear: 1988,
      },
      borrowedAt: "2026-07-08T10:00:00Z",
      state: {
        status: "returned-late",
        returnedAt: "2026-07-18T10:00:00Z",
        overdueDays: 3,
        fine: { amountMinor: 15000, currency: "IDR", status: "unpaid" },
      },
    },
  ],
  summary: {
    activeLoans: 1,
    completedLoans: 1,
    lateReturns: 1,
    unpaidFineMinor: 15000,
    fineCurrency: "IDR",
  },
};

describe("DashboardClient", () => {
  beforeEach(() => {
    window.history.replaceState(null, "", "/dashboard?tab=books");
  });

  it("greets the authenticated member", () => {
    render(
      <DashboardClient
        library={{ status: "success", library }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    expect(screen.getByRole("heading", { name: /welcome, maya/i })).toBeVisible();
    expect(screen.getByRole("tab", { name: "My Books" })).toHaveAttribute(
      "aria-selected",
      "true",
    );
    expect(screen.getByRole("heading", { name: "My Books" })).toBeVisible();
    expect(screen.getByText("My History")).not.toBeVisible();
    expect(screen.getAllByText("Atomic Habits")).toHaveLength(2);
    expect(screen.getAllByText(/15[.,]000/).length).toBeGreaterThan(0);
  });

  it("switches tabs instantly and records the selection in browser history", async () => {
    const user = userEvent.setup();
    window.history.replaceState(null, "", "/dashboard?tab=books&source=borrow");
    render(
      <DashboardClient
        library={{ status: "success", library }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    await user.click(screen.getByRole("tab", { name: "History" }));

    expect(window.location.search).toBe("?tab=history&source=borrow");
    expect(screen.getByRole("heading", { name: "My History" })).toBeVisible();
    expect(screen.getByText("The Alchemist")).toBeVisible();
    expect(screen.getByText("3 days late")).toBeVisible();
    expect(screen.getByText("Unpaid")).toBeVisible();
  });

  it("restores the visible panel when browser history changes", async () => {
    const user = userEvent.setup();
    render(
      <DashboardClient
        library={{ status: "success", library }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    await user.click(screen.getByRole("tab", { name: "History" }));
    act(() => {
      window.history.replaceState(null, "", "/dashboard?tab=books");
      window.dispatchEvent(new PopStateEvent("popstate"));
    });

    expect(screen.getByRole("heading", { name: "My Books" })).toBeVisible();
    expect(screen.getByText("My History")).not.toBeVisible();
  });

  it("moves tab focus with arrow keys without changing the active panel", async () => {
    const user = userEvent.setup();
    render(
      <DashboardClient
        library={{ status: "success", library }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    const booksTab = screen.getByRole("tab", { name: "My Books" });
    const historyTab = screen.getByRole("tab", { name: "History" });
    booksTab.focus();
    await user.keyboard("{ArrowRight}");

    expect(historyTab).toHaveFocus();
    expect(booksTab).toHaveAttribute("aria-selected", "true");
  });

  it("posts logout through the server session route", () => {
    render(
      <DashboardClient
        library={{ status: "success", library }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    const form = screen.getByRole("button", { name: "Log out" }).closest("form");
    expect(form).toHaveAttribute("action", "/api/auth/logout");
    expect(form).toHaveAttribute("method", "post");
  });

  it("shows an actionable error without hiding the member account", () => {
    render(
      <DashboardClient
        library={{ status: "error", error: { kind: "unavailable" } }}
        initialTab="books"
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        notice={{ kind: "none" }}
        session={session}
      />,
    );

    expect(screen.getByRole("heading", { name: /welcome, maya chen/i })).toBeVisible();
    expect(screen.getByRole("alert")).toHaveTextContent(
      "Library activity could not be loaded",
    );
  });
});
