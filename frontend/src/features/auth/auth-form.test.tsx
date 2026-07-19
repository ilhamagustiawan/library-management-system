import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, expect, it, vi } from "vitest";

import { LoginForm } from "./login-form";
import { RegisterForm } from "./register-form";

function renderWithQueryClient(component: React.ReactNode) {
  const queryClient = new QueryClient({
    defaultOptions: { mutations: { retry: false } },
  });

  return render(
    <QueryClientProvider client={queryClient}>{component}</QueryClientProvider>,
  );
}

describe("LoginForm", () => {
  it("shows actionable validation errors", async () => {
    const user = userEvent.setup();
    renderWithQueryClient(<LoginForm onAuthenticated={vi.fn()} />);

    await user.type(screen.getByLabelText("Email address"), "wrong");
    await user.type(screen.getByLabelText("Password"), "short");
    await user.click(screen.getByRole("button", { name: "Log in" }));

    expect(await screen.findByText("Enter a valid email address.")).toBeVisible();
    expect(screen.getByText("Password must contain at least 8 characters.")).toBeVisible();
  });

  it("returns a mock session for valid details", async () => {
    const user = userEvent.setup();
    const onAuthenticated = vi.fn();
    renderWithQueryClient(<LoginForm onAuthenticated={onAuthenticated} />);

    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password"), "quietreading");
    await user.click(screen.getByRole("button", { name: "Log in" }));

    await waitFor(() => expect(onAuthenticated).toHaveBeenCalledOnce());
    expect(onAuthenticated).toHaveBeenCalledWith(
      expect.objectContaining({ email: "maya@libry.test" }),
    );
  });
});

describe("RegisterForm", () => {
  it("requires matching passwords and terms acceptance", async () => {
    const user = userEvent.setup();
    renderWithQueryClient(<RegisterForm onAuthenticated={vi.fn()} />);

    await user.type(screen.getByLabelText("Full name"), "Maya Chen");
    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password", { selector: "#password" }), "quietreading");
    await user.type(screen.getByLabelText("Confirm password"), "differentpassword");
    await user.click(screen.getByRole("button", { name: "Create account" }));

    expect(await screen.findByText("Passwords must match.")).toBeVisible();
    expect(screen.getByText("Accept the terms to create an account.")).toBeVisible();
  });

  it("returns a mock session for valid details", async () => {
    const user = userEvent.setup();
    const onAuthenticated = vi.fn();
    renderWithQueryClient(<RegisterForm onAuthenticated={onAuthenticated} />);

    await user.type(screen.getByLabelText("Full name"), "Maya Chen");
    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password", { selector: "#password" }), "quietreading");
    await user.type(screen.getByLabelText("Confirm password"), "quietreading");
    await user.click(screen.getByRole("checkbox"));
    await user.click(screen.getByRole("button", { name: "Create account" }));

    await waitFor(() => expect(onAuthenticated).toHaveBeenCalledOnce());
    expect(onAuthenticated).toHaveBeenCalledWith({
      id: "mock-member",
      name: "Maya Chen",
      email: "maya@libry.test",
    });
  });
});
