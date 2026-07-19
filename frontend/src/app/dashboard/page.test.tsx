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
    redirect: vi.fn((target: string): never => {
      throw new Error(`redirect:${target}`);
    }),
  };
});

vi.mock("next/headers", () => ({
  cookies: async () => ({ get: () => ({ value: next.cookieValue }) }),
}));
vi.mock("next/navigation", () => ({ redirect: next.redirect }));
vi.mock("@/features/auth/session-refresh", () => ({
  SessionRefresh: () => <p>Refreshing secure session…</p>,
}));

const sessionSecret = "abcdef0123456789abcdef0123456789";

describe("DashboardPage", () => {
  beforeEach(() => {
    next.cookieValue = undefined;
    next.redirect.mockClear();
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8081");
    vi.stubEnv("AUTH_CLIENT_ID", "nextjs");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", "http://localhost:3000/api/auth/callback/library");
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
  });

  it("sends guests through the login page browser-navigation boundary", async () => {
    await expect(DashboardPage()).rejects.toThrow("redirect:/login");
  });

  it("renders a POST refresh boundary for an expiring session", async () => {
    next.cookieValue = WebSession.seal(
      {
        user: { id: "user-123", name: "Maya Chen", email: "maya@libry.test" },
        accessToken: "access-token",
        refreshToken: "refresh-token",
        tokenType: "Bearer",
        scope: "library:read library:write",
        expiresAt: 1,
      },
      sessionSecret,
    );

    render(await DashboardPage());

    expect(screen.getByText("Refreshing secure session…")).toBeVisible();
    expect(next.redirect).not.toHaveBeenCalledWith("/api/auth/refresh?return_to=/dashboard");
  });
});
