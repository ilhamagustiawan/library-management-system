import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { AuthSession } from "./auth-session";
import { DashboardClient } from "./dashboard-client";

const navigation = vi.hoisted(() => ({
  replace: vi.fn(),
}));

vi.mock("next/navigation", () => ({
  useRouter: () => navigation,
}));

describe("DashboardClient", () => {
  beforeEach(() => {
    window.localStorage.clear();
    navigation.replace.mockClear();
  });

  it("redirects guests to login", async () => {
    render(<DashboardClient />);

    await waitFor(() => expect(navigation.replace).toHaveBeenCalledWith("/login"));
  });

  it("greets the stored member", async () => {
    AuthSession.write(window.localStorage, {
      id: "mock-member",
      name: "Maya Chen",
      email: "maya@libry.test",
    });

    render(<DashboardClient />);

    expect(await screen.findByRole("heading", { name: /welcome, maya/i })).toBeVisible();
  });

  it("clears the session on logout", async () => {
    const user = userEvent.setup();
    AuthSession.write(window.localStorage, {
      id: "mock-member",
      name: "Maya Chen",
      email: "maya@libry.test",
    });
    render(<DashboardClient />);

    await user.click(await screen.findByRole("button", { name: "Log out" }));

    expect(AuthSession.read(window.localStorage)).toBeNull();
    expect(navigation.replace).toHaveBeenCalledWith("/");
  });
});
