import { describe, expect, it, vi } from "vitest";

import { MemberLibrary } from "./member-library";

const issuer = "http://localhost:8000";

function transactionPage() {
  return Response.json({
    code: "LMS-200000",
    data: {
      items: [
        {
          id: "transaction-return-late",
          loanId: "loan-late",
          memberId: "member-1",
          bookId: "book-late",
          type: "return",
          occurredAt: "2026-07-18T10:00:00Z",
          fine: {
            id: "fine-1",
            loanId: "loan-late",
            memberId: "member-1",
            overdueDays: 3,
            dailyRateMinor: 5000,
            totalAmountMinor: 15000,
            currency: "IDR",
            status: "unpaid",
            assessedAt: "2026-07-18T10:00:00Z",
          },
        },
        {
          id: "transaction-borrow-active",
          loanId: "loan-active",
          memberId: "member-1",
          bookId: "book-active",
          type: "borrow",
          occurredAt: "2026-07-17T10:00:00Z",
        },
        {
          id: "transaction-borrow-late",
          loanId: "loan-late",
          memberId: "member-1",
          bookId: "book-late",
          type: "borrow",
          occurredAt: "2026-07-08T10:00:00Z",
        },
        {
          id: "transaction-return-on-time",
          loanId: "loan-on-time",
          memberId: "member-1",
          bookId: "book-on-time",
          type: "return",
          occurredAt: "2026-07-05T10:00:00Z",
        },
        {
          id: "transaction-borrow-on-time",
          loanId: "loan-on-time",
          memberId: "member-1",
          bookId: "book-on-time",
          type: "borrow",
          occurredAt: "2026-07-01T10:00:00Z",
        },
      ],
      page: 1,
      pageSize: 100,
      totalItems: 5,
      totalPages: 1,
    },
  });
}

function bookResponse(id: string) {
  const books: Record<string, { title: string; author: string }> = {
    "book-active": { title: "Atomic Habits", author: "James Clear" },
    "book-late": { title: "The Alchemist", author: "Paulo Coelho" },
    "book-on-time": { title: "Clean Code", author: "Robert C. Martin" },
  };
  const book = books[id];
  if (book === undefined) return Response.json({}, { status: 404 });
  return Response.json({
    code: "LMS-200000",
    data: {
      id,
      isbn: "9780132350884",
      title: book.title,
      author: book.author,
      description: null,
      coverUrl: null,
      publicationYear: 2008,
      totalCopies: 3,
      availableCopies: 2,
      createdAt: "2026-07-01T00:00:00Z",
      updatedAt: "2026-07-01T00:00:00Z",
    },
  });
}

describe("MemberLibrary.load", () => {
  it("builds active books and loan history with late-return fine details", async () => {
    const fetcher = vi.fn<typeof fetch>(async (input) => {
      const url = new URL(input.toString());
      if (url.pathname === "/api/v1/transactions/me") return transactionPage();
      return bookResponse(url.pathname.split("/").at(-1) ?? "");
    });

    const result = await MemberLibrary.load({ issuer, accessToken: "access-token", fetcher });

    expect(result.status).toBe("success");
    if (result.status !== "success") return;
    expect(result.library.activeLoans).toEqual([
      expect.objectContaining({ loanId: "loan-active", book: expect.objectContaining({ title: "Atomic Habits" }) }),
    ]);
    expect(result.library.history).toEqual([
      expect.objectContaining({
        loanId: "loan-late",
        state: {
          status: "returned-late",
          returnedAt: "2026-07-18T10:00:00Z",
          overdueDays: 3,
          fine: { amountMinor: 15000, currency: "IDR", status: "unpaid" },
        },
      }),
      expect.objectContaining({ loanId: "loan-active", state: { status: "active" } }),
      expect.objectContaining({
        loanId: "loan-on-time",
        state: { status: "returned-on-time", returnedAt: "2026-07-05T10:00:00Z" },
      }),
    ]);
    expect(result.library.summary).toEqual({
      activeLoans: 1,
      completedLoans: 2,
      lateReturns: 1,
      unpaidFineMinor: 15000,
      fineCurrency: "IDR",
    });
    expect(MemberLibrary.activeLoanForBook(result, "book-active")?.loanId).toBe("loan-active");
    expect(MemberLibrary.activeLoanForBook(result, "book-late")).toBeUndefined();
  });

  it("returns a typed error when transaction history is unavailable", async () => {
    const result = await MemberLibrary.load({
      issuer,
      accessToken: "access-token",
      fetcher: vi.fn<typeof fetch>(async () => new Response(null, { status: 503 })),
    });

    expect(result).toEqual({ status: "error", error: { kind: "unavailable" } });
  });

  it("keeps every active loan reachable from History", async () => {
    const completed = Array.from({ length: 21 }, (_, index) => {
      const day = String(30 - index).padStart(2, "0");
      return [
        {
          id: `return-${index}`,
          loanId: `completed-${index}`,
          memberId: "member-1",
          bookId: `book-${index}`,
          type: "return",
          occurredAt: `2026-06-${day}T10:00:00Z`,
        },
        {
          id: `borrow-${index}`,
          loanId: `completed-${index}`,
          memberId: "member-1",
          bookId: `book-${index}`,
          type: "borrow",
          occurredAt: `2026-05-${day}T10:00:00Z`,
        },
      ];
    }).flat();
    const fetcher = vi.fn<typeof fetch>(async (input) => {
      const url = new URL(input.toString());
      if (url.pathname === "/api/v1/transactions/me") {
        return Response.json({
          code: "LMS-200000",
          data: {
            items: [
              ...completed,
              {
                id: "borrow-active-old",
                loanId: "active-old",
                memberId: "member-1",
                bookId: "book-active-old",
                type: "borrow",
                occurredAt: "2026-01-01T10:00:00Z",
              },
            ],
            page: 1,
            pageSize: 100,
            totalItems: 43,
            totalPages: 1,
          },
        });
      }
      return new Response(null, { status: 404 });
    });

    const result = await MemberLibrary.load({ issuer, accessToken: "access-token", fetcher });

    expect(result.status).toBe("success");
    if (result.status !== "success") return;
    expect(result.library.history.some((loan) => loan.loanId === "active-old")).toBe(true);
  });
});
