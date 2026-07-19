import type { CatalogBook } from "./book-catalog";
import type { LoadMemberLibraryResult } from "./member-library";

export type LoanEligibility =
  | { status: "eligible" }
  | {
      status: "disabled";
      reason: "loan-limit" | "already-borrowed" | "no-copies" | "account-unavailable";
      message: string;
    };

function forBook(book: CatalogBook, library: LoadMemberLibraryResult): LoanEligibility {
  if (book.availableCopies === 0) {
    return { status: "disabled", reason: "no-copies", message: "No copies currently available." };
  }
  if (library.status === "error") {
    return {
      status: "disabled",
      reason: "account-unavailable",
      message: "Loan status could not be checked. Try again before borrowing.",
    };
  }
  const alreadyBorrowed = library.library.activeLoans.some((loan) => loan.book.id === book.id);
  if (alreadyBorrowed) {
    return {
      status: "disabled",
      reason: "already-borrowed",
      message: "This book is already in your active loans.",
    };
  }
  if (library.library.summary.activeLoans >= 3) {
    return {
      status: "disabled",
      reason: "loan-limit",
      message: "Loan limit reached. Return a book before borrowing another.",
    };
  }
  return { status: "eligible" };
}

export const LoanEligibility = { forBook } as const;
