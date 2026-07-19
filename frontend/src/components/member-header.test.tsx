import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { MemberHeader } from "./member-header";

vi.mock("next/navigation", () => ({
  usePathname: () => "/dashboard",
}));

describe("MemberHeader", () => {
  it("returns authenticated members to their books from the title", () => {
    render(<MemberHeader logoutEndpoint="http://localhost:8000/api/v1/auth/logout" />);

    expect(screen.getByRole("link", { name: "Perpus Digital" })).toHaveAttribute(
      "href",
      "/dashboard?tab=books",
    );
  });

  it("marks dashboard navigation as current", () => {
    render(<MemberHeader logoutEndpoint="http://localhost:8000/api/v1/auth/logout" />);

    expect(screen.getByRole("link", { name: "My Books" })).toHaveAttribute(
      "aria-current",
      "page",
    );
  });
});
