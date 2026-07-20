import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import { LogoutForm } from "./logout-form";

describe("LogoutForm", () => {
  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it("submits local logout after auth-service logout", async () => {
    const user = userEvent.setup();
    const submit = vi.spyOn(HTMLFormElement.prototype, "submit").mockImplementation(() => undefined);
    vi.stubGlobal("fetch", vi.fn(async () => new Response(null, { status: 204 })));

    render(<LogoutForm logoutEndpoint="http://localhost:8000/api/v1/auth/logout" />);
    await user.click(screen.getByRole("button", { name: "Log out" }));

    await waitFor(() => expect(submit).toHaveBeenCalledOnce());
    const form = screen.getByRole("button", { name: "Logging out…" }).closest("form");
    expect(form).toHaveAttribute("action", "/api/auth/logout");
    expect(form).toHaveAttribute("method", "post");
  });

  it("submits local logout without waiting for auth-service logout", async () => {
    const user = userEvent.setup();
    const submit = vi.spyOn(HTMLFormElement.prototype, "submit").mockImplementation(() => undefined);
    vi.stubGlobal("fetch", vi.fn(() => new Promise<Response>(() => undefined)));

    render(<LogoutForm logoutEndpoint="http://localhost:8000/api/v1/auth/logout" />);
    await user.click(screen.getByRole("button", { name: "Log out" }));

    await waitFor(() => expect(submit).toHaveBeenCalledOnce());
  });

  it("supports member headers on light surfaces", () => {
    render(
      <LogoutForm
        logoutEndpoint="http://localhost:8000/api/v1/auth/logout"
        tone="light"
      />,
    );

    expect(screen.getByRole("button", { name: "Log out" })).toHaveClass("text-foreground");
  });
});
