import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { DashboardClient } from "./dashboard-client";

const session = {
  id: "user-123",
  name: "Maya Chen",
  email: "maya@libry.test",
};

describe("DashboardClient", () => {
  it("greets the authenticated member", () => {
    render(<DashboardClient logoutEndpoint="http://localhost:8081/api/v1/auth/logout" session={session} />);

    expect(screen.getByRole("heading", { name: /welcome, maya/i })).toBeVisible();
  });

  it("posts logout through the server session route", () => {
    render(<DashboardClient logoutEndpoint="http://localhost:8081/api/v1/auth/logout" session={session} />);

    const form = screen.getByRole("button", { name: "Log out" }).closest("form");
    expect(form).toHaveAttribute("action", "/api/auth/logout");
    expect(form).toHaveAttribute("method", "post");
  });
});
