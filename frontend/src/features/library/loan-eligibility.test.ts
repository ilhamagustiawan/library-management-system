import { describe, expect, it } from "vitest";

import type { CatalogBook } from "./book-catalog";
import { LoanEligibility } from "./loan-eligibility";
import type { MemberLibrary } from "./member-library";

const book: CatalogBook = {
  id: "0ec82798-8ff9-48c5-b68f-2b8c050647ac",
  isbn: "9780132350884",
  title: "Clean Code",
  author: "Robert C. Martin",
  description: null,
  coverUrl: null,
  publicationYear: 2008,
  totalCopies: 3,
  availableCopies: 2,
  createdAt: "2026-07-01T00:00:00Z",
  updatedAt: "2026-07-01T00:00:00Z",
};

function library(activeLoans: number): MemberLibrary {
  return {
    activeLoans: [],
    history: [],
    summary: {
      activeLoans,
      completedLoans: 0,
      lateReturns: 0,
      unpaidFineMinor: 0,
      fineCurrency: "IDR",
    },
  };
}

describe("LoanEligibility", () => {
  it("disables borrowing at the three-loan API limit", () => {
    expect(LoanEligibility.forBook(book, { status: "success", library: library(3) })).toEqual({
      status: "disabled",
      reason: "loan-limit",
      message: "Loan limit reached. Return a book before borrowing another.",
    });
  });

  it("allows borrowing below the limit when stock exists", () => {
    expect(LoanEligibility.forBook(book, { status: "success", library: library(2) })).toEqual({
      status: "eligible",
    });
  });
});
