import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import { BorrowButton } from "./borrow-button";

const navigation = vi.hoisted(() => ({ push: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ push: navigation.push }),
}));

describe("BorrowButton", () => {
  afterEach(() => {
    navigation.push.mockReset();
    vi.unstubAllGlobals();
  });

  it("navigates with a one-time success marker after borrowing", async () => {
    const user = userEvent.setup();
    vi.stubGlobal(
      "fetch",
      vi.fn<typeof fetch>(async () => Response.json({ status: "borrowed" }, { status: 201 })),
    );
    render(
      <BorrowButton
        bookId="0ec82798-8ff9-48c5-b68f-2b8c050647ac"
        eligibility={{ status: "eligible" }}
      />,
    );

    await user.click(screen.getByRole("button", { name: "Borrow this book" }));

    await waitFor(() => expect(navigation.push).toHaveBeenCalledWith("/dashboard?borrowed=1"));
  });
});
