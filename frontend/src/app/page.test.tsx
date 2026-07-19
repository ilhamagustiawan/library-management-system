import { render, screen, within } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { WebSession } from "@/features/auth/web-session";
import HomePage from "./page";

const next = vi.hoisted(() => {
  function initialCookieValue(): string | undefined {
    return undefined;
  }
  return { cookieValue: initialCookieValue() };
});

vi.mock("next/headers", () => ({
  cookies: async () => ({
    get: () =>
      next.cookieValue === undefined ? undefined : { value: next.cookieValue },
  }),
}));

const sessionSecret = "abcdef0123456789abcdef0123456789";

describe("HomePage", () => {
  beforeEach(() => {
    next.cookieValue = undefined;
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8000");
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", "http://localhost:3000/api/auth/callback/library");
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
  });

  it("offers guests clear login and registration paths", async () => {
    render(await HomePage());

    expect(
      within(screen.getByRole("banner")).getByRole("link", { name: "Perpus Digital" }),
    ).toHaveAttribute("href", "/");
    expect(screen.getByRole("navigation", { name: "Primary" })).toBeVisible();
    expect(
      screen.getByRole("heading", { name: "Find your next read. Keep every loan in view." }),
    ).toBeVisible();
    expect(
      screen.getByRole("img", { name: /open book on a library reading desk/i }),
    ).toBeVisible();
    expect(screen.getByText(/browse books currently available/i)).toBeVisible();
    expect(screen.getByRole("link", { name: "Create my member account" })).toHaveAttribute(
      "href",
      "/register",
    );
    expect(screen.getByRole("link", { name: "Log in to my account" })).toHaveAttribute(
      "href",
      "/login",
    );
    expect(
      within(screen.getByRole("navigation", { name: "Reading desk" })).getByRole("link", {
        name: "Browse books",
      }),
    ).toHaveAttribute("href", "/browse");
  });

  it("shows member navigation instead of another login after authentication", async () => {
    next.cookieValue = WebSession.seal(
      {
        user: { id: "user-123", name: "Maya Chen", email: "maya@perpus-digital.test" },
        accessToken: "access-token",
        refreshToken: "refresh-token",
        tokenType: "Bearer",
        scope: "books:read loans:borrow:self",
        expiresAt: Math.floor(Date.now() / 1_000) + 900,
      },
      sessionSecret,
    );

    render(await HomePage());

    expect(within(screen.getByRole("banner")).getByRole("link", { name: "Perpus Digital" })).toHaveAttribute(
      "href",
      "/dashboard?tab=books",
    );
    expect(screen.getByRole("link", { name: "Member home" })).toHaveAttribute(
      "href",
      "/dashboard",
    );
    expect(screen.getByRole("button", { name: "Log out" })).toBeVisible();
    expect(
      within(screen.getByRole("navigation", { name: "Primary" })).queryByRole("link", {
        name: "Log in",
      }),
    ).not.toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Open member home" })).toHaveAttribute(
      "href",
      "/dashboard",
    );
  });
});
