import { render, screen, waitFor } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { OAuthStart } from "./oauth-start";

describe("OAuthStart", () => {
  it("starts OAuth with a full browser navigation", async () => {
    const navigate = vi.fn();

    render(<OAuthStart navigate={navigate} />);

    await waitFor(() => expect(navigate).toHaveBeenCalledWith("/api/auth/login"));
    expect(screen.getByRole("link", { name: "Continue to login" })).toHaveAttribute(
      "href",
      "/api/auth/login",
    );
  });
});
