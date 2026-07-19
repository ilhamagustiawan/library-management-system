import { beforeEach, describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";

import { WebSession } from "@/features/auth/web-session";
import DashboardPage from "./page";

const next = vi.hoisted(() => {
  function initialCookieValue(): string | undefined {
    return undefined;
  }
  return {
    cookieValue: initialCookieValue(),
    loadLibrary: vi.fn(),
    redirect: vi.fn((target: string): never => {
      throw new Error(`redirect:${target}`);
    }),
  };
});

vi.mock("next/headers", () => ({
  cookies: async () => ({ get: () => ({ value: next.cookieValue }) }),
}));
vi.mock("next/navigation", () => ({
  redirect: next.redirect,
  usePathname: () => "/dashboard",
}));
vi.mock("@/features/library/member-library", () => ({
  MemberLibrary: { load: next.loadLibrary },
}));
vi.mock("@/features/auth/session-refresh", () => ({
  SessionRefresh: () => <p>Refreshing secure session…</p>,
}));

const sessionSecret = "abcdef0123456789abcdef0123456789";

describe("DashboardPage", () => {
  beforeEach(() => {
    next.cookieValue = undefined;
    next.redirect.mockClear();
    next.loadLibrary.mockReset();
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8000");
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", "http://localhost:3000/api/auth/callback/library");
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
  });

  it("sends guests through the login page browser-navigation boundary", async () => {
    await expect(DashboardPage({ searchParams: Promise.resolve({}) })).rejects.toThrow(
      "redirect:/login",
    );
  });

  it("renders a POST refresh boundary for an expiring session", async () => {
    next.cookieValue = WebSession.seal(
      {
        user: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
        accessToken: "access-token",
        refreshToken: "refresh-token",
        tokenType: "Bearer",
        scope: "books:read loans:borrow:self",
        expiresAt: 1,
      },
      sessionSecret,
    );

    render(await DashboardPage({ searchParams: Promise.resolve({}) }));

    expect(screen.getByText("Refreshing secure session…")).toBeVisible();
    expect(next.redirect).not.toHaveBeenCalledWith("/api/auth/refresh?return_to=/dashboard");
  });

  it("loads member books and history with the server-held access token", async () => {
    next.cookieValue = WebSession.seal(
      {
        user: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
        accessToken: "access-token",
        refreshToken: "refresh-token",
        tokenType: "Bearer",
        scope: "books:read transactions:read:self",
        expiresAt: Math.floor(Date.now() / 1_000) + 900,
      },
      sessionSecret,
    );
    next.loadLibrary.mockResolvedValue({
      status: "success",
      library: {
        activeLoans: [],
        history: [],
        summary: {
          activeLoans: 0,
          completedLoans: 0,
          lateReturns: 0,
          unpaidFineMinor: 0,
          fineCurrency: "IDR",
        },
      },
    });

    render(await DashboardPage({ searchParams: Promise.resolve({}) }));

    expect(screen.getByRole("heading", { name: "My Books" })).toBeVisible();
    expect(screen.getByText("My History")).not.toBeVisible();
    expect(next.loadLibrary).toHaveBeenCalledWith(
      expect.objectContaining({ issuer: "http://localhost:8000", accessToken: "access-token" }),
    );
  });
});
