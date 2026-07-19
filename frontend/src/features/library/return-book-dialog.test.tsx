import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import { Button } from "@/components/ui/button";

import { ReturnBookDialog } from "./return-book-dialog";

const navigation = vi.hoisted(() => ({ refresh: vi.fn() }));
const notifications = vi.hoisted(() => ({ success: vi.fn() }));

vi.mock("next/navigation", () => ({
  useRouter: () => ({ refresh: navigation.refresh }),
}));

vi.mock("sonner", () => ({
  toast: { success: notifications.success },
}));

const loanId = "52a88672-a4c2-4876-be5a-65863aeb35e4";

function quote(totalAmountMinor: number) {
  return {
    status: "ready",
    quote: {
      loanId,
      bookId: "7b36fe43-f31d-4861-884f-42ed7386b1e9",
      dueAt: "2026-07-26T10:00:00Z",
      quotedAt: "2026-07-28T10:00:00Z",
      fine:
        totalAmountMinor === 0
          ? null
          : {
              overdueDays: 2,
              dailyRateMinor: 5000,
              totalAmountMinor,
              currency: "IDR",
            },
    },
  };
}

function renderDialog() {
  render(
    <ReturnBookDialog bookTitle="Clean Code" loanId={loanId}>
      <Button>Return book</Button>
    </ReturnBookDialog>,
  );
}

describe("ReturnBookDialog", () => {
  afterEach(() => {
    navigation.refresh.mockReset();
    notifications.success.mockReset();
    vi.unstubAllGlobals();
  });

  it("warns about the exact late fine before returning", async () => {
    const user = userEvent.setup();
    const fetcher = vi.fn<typeof fetch>(async (_input, options) => {
      if (options?.method === "POST") {
        return Response.json(
          {
            status: "returned",
            stockUpdate: "pending",
            fine: { overdueDays: 2, totalAmountMinor: 10000, currency: "IDR" },
          },
          { status: 202 },
        );
      }
      return Response.json(quote(10000));
    });
    vi.stubGlobal("fetch", fetcher);
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Return book" }));
    const dialog = await screen.findByRole("dialog");

    expect(within(dialog).getByText("2 days late")).toBeInTheDocument();
    expect(within(dialog).getByText(/Rp\s*10\.000/)).toBeInTheDocument();
    await user.click(within(dialog).getByRole("button", { name: "Confirm return" }));

    await waitFor(() => expect(navigation.refresh).toHaveBeenCalledOnce());
    const post = fetcher.mock.calls.find((call) => call[1]?.method === "POST");
    expect(post?.[1]?.body).toBe(JSON.stringify({ acceptedFineAmountMinor: 10000 }));
    expect(notifications.success).toHaveBeenCalledWith(
      "Book returned. Fine assessed: Rp 10.000",
      expect.objectContaining({ description: "Return recorded; catalog availability is updating." }),
    );
  });

  it("requires confirmation for an on-time return", async () => {
    const user = userEvent.setup();
    vi.stubGlobal("fetch", vi.fn<typeof fetch>(async () => Response.json(quote(0))));
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Return book" }));

    const dialog = await screen.findByRole("dialog");
    expect(within(dialog).getByText("No fine will be assessed.")).toBeInTheDocument();
    expect(within(dialog).getByRole("button", { name: "Confirm return" })).toBeEnabled();
  });

  it("reloads a changed fine and requires another confirmation", async () => {
    const user = userEvent.setup();
    let quoteRequests = 0;
    const fetcher = vi.fn<typeof fetch>(async (_input, options) => {
      if (options?.method === "POST") {
        return Response.json(
          {
            error: {
              kind: "fine-quote-changed",
              message: "Fine changed before return. Review the updated amount and confirm again.",
            },
          },
          { status: 409 },
        );
      }
      quoteRequests += 1;
      return Response.json(quote(quoteRequests === 1 ? 0 : 5000));
    });
    vi.stubGlobal("fetch", fetcher);
    renderDialog();

    await user.click(screen.getByRole("button", { name: "Return book" }));
    await user.click(await screen.findByRole("button", { name: "Confirm return" }));

    expect(await screen.findByText("Fine changed before return. Review the updated amount and confirm again.")).toBeInTheDocument();
    expect(screen.getAllByText(/Rp\s*5\.000/)).not.toHaveLength(0);
    expect(navigation.refresh).not.toHaveBeenCalled();
  });
});
