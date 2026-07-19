import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

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
  const loginEndpoint = "http://localhost:8081/api/v1/auth/login";
  const returnTo = "http://localhost:8081/oauth/authorize?client_id=nextjs";

  afterEach(() => vi.unstubAllGlobals());

  it("shows actionable validation errors", async () => {
    const user = userEvent.setup();
    renderWithQueryClient(<LoginForm loginEndpoint={loginEndpoint} returnTo={returnTo} />);

    await user.type(screen.getByLabelText("Email address"), "wrong");
    await user.type(screen.getByLabelText("Password"), "short");
    await user.click(screen.getByRole("button", { name: "Log in" }));

    expect(await screen.findByText("Enter a valid email address.")).toBeVisible();
  });

  it("continues authorization after auth service login", async () => {
    const user = userEvent.setup();
    const navigate = vi.fn();
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json({
          code: "LMS-200000",
          data: { id: "user-123", name: "Maya Chen", email: "maya@libry.test" },
        }),
      ),
    );
    renderWithQueryClient(
      <LoginForm loginEndpoint={loginEndpoint} navigate={navigate} returnTo={returnTo} />,
    );

    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password"), "quietreading");
    await user.click(screen.getByRole("button", { name: "Log in" }));

    await waitFor(() => expect(navigate).toHaveBeenCalledWith(returnTo));
  });
});

describe("RegisterForm", () => {
  const registerEndpoint = "http://localhost:8081/api/v1/auth/register";

  it("requires matching passwords and terms acceptance", async () => {
    const user = userEvent.setup();
    renderWithQueryClient(<RegisterForm registerEndpoint={registerEndpoint} />);

    await user.type(screen.getByLabelText("Full name"), "Maya Chen");
    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password", { selector: "#password" }), "quietreading");
    await user.type(screen.getByLabelText("Confirm password"), "differentpassword");
    await user.click(screen.getByRole("button", { name: "Create account" }));

    expect(await screen.findByText("Passwords must match.")).toBeVisible();
    expect(screen.getByText("Accept the terms to create an account.")).toBeVisible();
  });

  it("continues to login after real registration", async () => {
    const user = userEvent.setup();
    const navigate = vi.fn();
    vi.stubGlobal(
      "fetch",
      vi.fn(async () =>
        Response.json(
          {
            code: "LMS-200000",
            data: { id: "user-123", name: "Maya Chen", email: "maya@libry.test" },
          },
          { status: 201 },
        ),
      ),
    );
    renderWithQueryClient(
      <RegisterForm registerEndpoint={registerEndpoint} navigate={navigate} />,
    );

    await user.type(screen.getByLabelText("Full name"), "Maya Chen");
    await user.type(screen.getByLabelText("Email address"), "maya@libry.test");
    await user.type(screen.getByLabelText("Password", { selector: "#password" }), "quietreading");
    await user.type(screen.getByLabelText("Confirm password"), "quietreading");
    await user.click(screen.getByRole("checkbox"));
    await user.click(screen.getByRole("button", { name: "Create account" }));

    await waitFor(() => expect(navigate).toHaveBeenCalledWith("/login"));
  });
});
