import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import type { MemberLoan } from "./member-library";
import { MemberBooksPanel, MemberHistoryPanel } from "./member-loan-panels";

vi.mock("next/navigation", () => ({ useRouter: () => ({ refresh: vi.fn() }) }));

const activeLoan: MemberLoan = {
  loanId: "52a88672-a4c2-4876-be5a-65863aeb35e4",
  borrowedAt: "2026-07-19T10:00:00Z",
  state: { status: "active" },
  book: {
    status: "available",
    id: "7b36fe43-f31d-4861-884f-42ed7386b1e9",
    title: "Clean Code",
    author: "Robert C. Martin",
    coverUrl: null,
    publicationYear: 2008,
  },
};

const returnedLoan: MemberLoan = {
  ...activeLoan,
  loanId: "16ffadc2-9c66-410d-90a9-29f7ff8399a1",
  state: { status: "returned-on-time", returnedAt: "2026-07-25T10:00:00Z" },
};

describe("member loan panels", () => {
  it("offers return from My Books and links to book detail", () => {
    render(<MemberBooksPanel loans={[activeLoan]} />);

    expect(screen.getByRole("button", { name: "Return book" })).toBeInTheDocument();
    expect(screen.getByRole("link", { name: "Clean Code" })).toHaveAttribute(
      "href",
      `/books/${activeLoan.book.id}`,
    );
  });

  it("offers return only for active History rows", () => {
    render(<MemberHistoryPanel loans={[activeLoan, returnedLoan]} />);

    expect(screen.getAllByRole("button", { name: "Return book" })).toHaveLength(1);
  });
});
