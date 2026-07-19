import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import LoginPage from "./login/page";
import RegisterPage from "./register/page";

vi.mock("next/navigation", () => ({ useRouter: () => ({ replace: vi.fn() }) }));
vi.mock("@/features/auth/oauth-start", () => ({
  OAuthStart: () => <a href="/api/auth/login">Continue to login</a>,
}));

function renderPage(page: React.ReactNode) {
  return render(
    <QueryClientProvider client={new QueryClient()}>{page}</QueryClientProvider>,
  );
}

describe("authentication pages", () => {
  beforeEach(() => {
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8000");
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", "http://localhost:3000/api/auth/callback/library");
    vi.stubEnv("AUTH_SESSION_SECRET", "abcdef0123456789abcdef0123456789");
  });

  it("uses the login title as the page heading", async () => {
    const page = await LoginPage({
      searchParams: Promise.resolve({
        return_to: "http://localhost:8000/oauth/authorize?client_id=member-nextjs-web",
      }),
    });
    renderPage(page);

    expect(
      screen.getByRole("heading", { level: 1, name: "Open your member account" }),
    ).toBeVisible();
  });

  it("renders a browser OAuth start instead of a server redirect", async () => {
    const page = await LoginPage({ searchParams: Promise.resolve({}) });
    renderPage(page);

    expect(screen.getByRole("link", { name: "Continue to login" })).toHaveAttribute(
      "href",
      "/api/auth/login",
    );
  });

  it("uses the registration title as the page heading", () => {
    renderPage(<RegisterPage />);

    expect(
      screen.getByRole("heading", { level: 1, name: "Create your member account" }),
    ).toBeVisible();
  });
});
