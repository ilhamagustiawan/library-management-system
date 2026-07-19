import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import LoginPage from "./login/page";
import RegisterPage from "./register/page";

vi.mock("next/navigation", () => ({
  useRouter: () => ({ replace: vi.fn() }),
}));

function renderPage(page: React.ReactNode) {
  return render(
    <QueryClientProvider client={new QueryClient()}>{page}</QueryClientProvider>,
  );
}

describe("authentication pages", () => {
  it("uses the login title as the page heading", () => {
    renderPage(<LoginPage />);

    expect(
      screen.getByRole("heading", { level: 1, name: "Open your member account" }),
    ).toBeVisible();
  });

  it("uses the registration title as the page heading", () => {
    renderPage(<RegisterPage />);

    expect(
      screen.getByRole("heading", { level: 1, name: "Create your member account" }),
    ).toBeVisible();
  });
});
