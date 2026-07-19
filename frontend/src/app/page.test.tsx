import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import HomePage from "./page";

describe("HomePage", () => {
  it("offers clear login and registration paths", () => {
    render(<HomePage />);

    expect(screen.getByRole("link", { name: "Create member account" })).toHaveAttribute(
      "href",
      "/register",
    );
    expect(screen.getByRole("link", { name: "I already have an account" })).toHaveAttribute(
      "href",
      "/login",
    );
  });
});
