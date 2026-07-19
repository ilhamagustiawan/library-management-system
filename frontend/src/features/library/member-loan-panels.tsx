import { AlertTriangle, BookCheck } from "lucide-react";
import Link from "next/link";
import type { ReactElement } from "react";

import { Button } from "@/components/ui/button";

import { BookCover } from "./book-cover";
import type { BookReference, MemberLoan } from "./member-library";
import { ReturnBookDialog } from "./return-book-dialog";

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en", {
    day: "numeric",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  }).format(new Date(value));
}

function formatMoney(amountMinor: number, currency: "IDR") {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency,
    maximumFractionDigits: 0,
  }).format(amountMinor);
}

function BookIdentity({ book }: { book: BookReference }) {
  if (book.status === "unavailable") {
    return (
      <div>
        <p className="font-semibold text-card-foreground">Catalog record unavailable</p>
        <p className="mt-1 text-xs text-muted-foreground">Book ID {book.id}</p>
      </div>
    );
  }
  return (
    <div>
      <Link
        className="font-semibold leading-5 text-primary underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        href={`/books/${book.id}`}
      >
        {book.title}
      </Link>
      <p className="mt-1 text-xs leading-5 text-muted-foreground">
        {book.author}
        {book.publicationYear === null ? "" : ` · ${book.publicationYear}`}
      </p>
    </div>
  );
}

function LoanStatus({ loan }: { loan: MemberLoan }) {
  switch (loan.state.status) {
    case "active":
      return <span className="font-semibold text-primary">Checked out</span>;
    case "returned-on-time":
      return <span className="font-semibold text-primary">Returned on time</span>;
    case "returned-late":
      return (
        <span className="inline-flex items-center gap-1 font-semibold text-destructive">
          <AlertTriangle aria-hidden="true" className="size-3.5" />
          {loan.state.overdueDays} {loan.state.overdueDays === 1 ? "day" : "days"} late
        </span>
      );
  }
}

function LoanReturnAction({ loan, children }: { loan: MemberLoan; children: ReactElement }) {
  const title = loan.book.status === "available" ? loan.book.title : "Book";
  return <ReturnBookDialog bookTitle={title} loanId={loan.loanId}>{children}</ReturnBookDialog>;
}

export function MemberBooksPanel({ loans }: { loans: MemberLoan[] }) {
  return (
    <section
      aria-labelledby="my-books-heading"
      className="overflow-hidden rounded-b-lg border border-t-0 border-border bg-card"
    >
      <header className="flex items-end justify-between gap-4 border-b border-border bg-primary/5 px-4 py-4 sm:px-5">
        <div>
          <p className="text-xs font-semibold uppercase tracking-[0.12em] text-primary">
            Current shelf
          </p>
          <h2 id="my-books-heading" className="mt-1 font-display text-2xl font-semibold">
            My Books
          </h2>
        </div>
        <p className="text-xs text-muted-foreground">{loans.length} checked out</p>
      </header>
      {loans.length === 0 ? (
        <div className="px-5 py-10 text-center">
          <BookCheck aria-hidden="true" className="mx-auto size-8 text-primary" strokeWidth={1.5} />
          <p className="mt-3 font-semibold">No books checked out</p>
          <p className="mt-1 text-sm text-muted-foreground">New loans will appear here.</p>
          <Link
            className="mt-4 inline-flex h-9 items-center rounded-md border border-primary bg-primary px-4 text-sm font-semibold text-primary-foreground hover:bg-primary/90 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2"
            href="/browse"
          >
            Browse books
          </Link>
        </div>
      ) : (
        <ul className="divide-y divide-border">
          {loans.map((loan) => {
            const title = loan.book.status === "available" ? loan.book.title : "Book";
            const coverUrl = loan.book.status === "available" ? loan.book.coverUrl : null;
            return (
              <li
                key={loan.loanId}
                className="grid grid-cols-[5rem_minmax(0,1fr)] items-start gap-4 px-4 py-5 sm:grid-cols-[6rem_minmax(0,1fr)_9rem] sm:px-5"
              >
                <BookCover coverUrl={coverUrl} title={title} size="loan" />
                <BookIdentity book={loan.book} />
                <div className="col-start-2 text-xs leading-5 sm:col-start-auto sm:text-right">
                  <p className="text-muted-foreground">Borrowed {formatDate(loan.borrowedAt)}</p>
                  <p className="mt-1"><LoanStatus loan={loan} /></p>
                  <LoanReturnAction loan={loan}>
                    <Button className="mt-3 w-full sm:w-auto" size="sm" variant="outline">
                      Return book
                    </Button>
                  </LoanReturnAction>
                </div>
              </li>
            );
          })}
        </ul>
      )}
    </section>
  );
}

function FineDetails({ loan }: { loan: MemberLoan }) {
  if (loan.state.status !== "returned-late") {
    return <span className="text-muted-foreground">No fine</span>;
  }
  return (
    <span>
      <span className="block font-semibold text-destructive">
        {formatMoney(loan.state.fine.amountMinor, loan.state.fine.currency)}
      </span>
      <span className="mt-1 block text-xs font-semibold uppercase tracking-[0.08em] text-destructive">
        Unpaid
      </span>
    </span>
  );
}

export function MemberHistoryPanel({ loans }: { loans: MemberLoan[] }) {
  return (
    <section
      aria-labelledby="history-heading"
      className="overflow-hidden rounded-b-lg border border-t-0 border-border bg-card"
    >
      <header className="border-b border-border bg-book-rust/5 px-4 py-4 sm:px-5">
        <p className="text-xs font-semibold uppercase tracking-[0.12em] text-primary">Loan record</p>
        <h2 id="history-heading" className="mt-1 font-display text-2xl font-semibold">My History</h2>
      </header>
      {loans.length === 0 ? (
        <p className="px-5 py-10 text-center text-sm text-muted-foreground">
          No borrowing history yet.
        </p>
      ) : (
        <div>
          <div
            aria-hidden="true"
            className="hidden grid-cols-[minmax(13rem,1fr)_9rem_10rem_8rem] gap-4 border-b border-border px-5 py-2 text-xs font-semibold text-muted-foreground md:grid"
          >
            <span>Book</span><span>Borrowed</span><span>Return status</span><span>Fine / action</span>
          </div>
          <ol className="divide-y divide-border">
            {loans.map((loan) => (
              <li
                key={loan.loanId}
                className="grid gap-4 px-4 py-4 text-sm md:grid-cols-[minmax(13rem,1fr)_9rem_10rem_8rem] md:px-5"
              >
                <BookIdentity book={loan.book} />
                <div><span className="mr-2 text-xs font-semibold text-muted-foreground md:hidden">Borrowed</span>{formatDate(loan.borrowedAt)}</div>
                <div>
                  <span className="mr-2 text-xs font-semibold text-muted-foreground md:hidden">Status</span>
                  <LoanStatus loan={loan} />
                  {loan.state.status !== "active" && (
                    <span className="mt-1 block text-xs text-muted-foreground">
                      {formatDate(loan.state.returnedAt)}
                    </span>
                  )}
                </div>
                <div>
                  <span className="mr-2 text-xs font-semibold text-muted-foreground md:hidden">
                    {loan.state.status === "active" ? "Action" : "Fine"}
                  </span>
                  {loan.state.status === "active" ? (
                    <LoanReturnAction loan={loan}>
                      <Button size="sm" variant="outline">Return book</Button>
                    </LoanReturnAction>
                  ) : (
                    <FineDetails loan={loan} />
                  )}
                </div>
              </li>
            ))}
          </ol>
        </div>
      )}
    </section>
  );
}
