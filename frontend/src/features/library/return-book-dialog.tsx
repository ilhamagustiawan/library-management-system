"use client";

import { AlertTriangle, LoaderCircle, RotateCcw } from "lucide-react";
import { useRouter } from "next/navigation";
import type { ReactElement } from "react";
import { useState } from "react";
import { toast } from "sonner";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

import { ReturnBook, type ReturnQuote } from "./return-book";

type State =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "ready"; quote: ReturnQuote; notice?: string }
  | { status: "submitting"; quote: ReturnQuote }
  | { status: "error"; message: string };

function formatMoney(amountMinor: number) {
  return new Intl.NumberFormat("id-ID", {
    style: "currency",
    currency: "IDR",
    maximumFractionDigits: 0,
  }).format(amountMinor);
}

function formatDate(value: string) {
  return new Intl.DateTimeFormat("en", {
    day: "numeric",
    month: "short",
    year: "numeric",
    timeZone: "UTC",
  }).format(new Date(value));
}

export function ReturnBookDialog({
  bookTitle,
  loanId,
  children,
}: {
  bookTitle: string;
  loanId: string;
  children: ReactElement;
}) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [state, setState] = useState<State>({ status: "idle" });

  async function loadQuote(notice?: string) {
    setState({ status: "loading" });
    const result = await ReturnBook.loadQuote(loanId);
    if (result.status === "error") {
      setState({ status: "error", message: result.message });
      return;
    }
    setState(
      notice === undefined
        ? { status: "ready", quote: result.quote }
        : { status: "ready", quote: result.quote, notice },
    );
  }

  function changeOpen(nextOpen: boolean) {
    if (!nextOpen && state.status === "submitting") return;
    setOpen(nextOpen);
    if (nextOpen) void loadQuote();
    else setState({ status: "idle" });
  }

  async function returnBook(quote: ReturnQuote) {
    setState({ status: "submitting", quote });
    const submission = await ReturnBook.submit(loanId, quote);
    if (submission.status === "quote-changed") {
      await loadQuote(submission.message);
      return;
    }
    if (submission.status === "error") {
      setState({ status: "error", message: submission.message });
      return;
    }
    setOpen(false);
    setState({ status: "idle" });
    const title = submission.result.fine === null
      ? "Book returned."
      : `Book returned. Fine assessed: ${formatMoney(submission.result.fine.totalAmountMinor)}`;
    toast.success(title, {
      description: submission.result.stockUpdate === "pending"
        ? "Return recorded; catalog availability is updating."
        : "Your library and catalog are now up to date.",
    });
    router.refresh();
  }

  return (
    <Dialog open={open} onOpenChange={changeOpen}>
      <DialogTrigger asChild>{children}</DialogTrigger>
      <DialogContent closable={state.status !== "submitting"}>
        <DialogHeader>
          <DialogTitle>Return {bookTitle}?</DialogTitle>
          <DialogDescription>
            Review the return details. Confirming records this loan as returned.
          </DialogDescription>
        </DialogHeader>

        {state.status === "loading" && (
          <p role="status" className="flex items-center gap-2 py-6 text-sm text-muted-foreground">
            <LoaderCircle aria-hidden="true" className="size-4 animate-spin motion-reduce:animate-none" />
            Checking due date and fine…
          </p>
        )}
        {state.status === "error" && (
          <Alert>
            <p className="font-semibold text-destructive">Return unavailable</p>
            <p className="mt-1 text-muted-foreground">{state.message}</p>
            <Button className="mt-3" size="sm" variant="outline" onClick={() => void loadQuote()}>
              <RotateCcw aria-hidden="true" /> Retry
            </Button>
          </Alert>
        )}
        {(state.status === "ready" || state.status === "submitting") && (
          <div className="space-y-4">
            {state.status === "ready" && state.notice !== undefined && (
              <Alert>
                <p className="font-semibold text-destructive">Fine updated</p>
                <p className="mt-1 text-muted-foreground">{state.notice}</p>
              </Alert>
            )}
            <p className="text-sm text-muted-foreground">
              Due <span className="font-semibold text-foreground">{formatDate(state.quote.dueAt)}</span>
            </p>
            {state.quote.fine === null ? (
              <div className="rounded-sm border border-primary/25 bg-primary/5 p-4 text-sm font-semibold text-primary">
                No fine will be assessed.
              </div>
            ) : (
              <Alert className="p-4">
                <p className="flex items-center gap-2 font-semibold text-destructive">
                  <AlertTriangle aria-hidden="true" className="size-4" />
                  {state.quote.fine.overdueDays} {state.quote.fine.overdueDays === 1 ? "day" : "days"} late
                </p>
                <p className="mt-2 text-2xl font-semibold text-destructive">
                  {formatMoney(state.quote.fine.totalAmountMinor)}
                </p>
                <p className="mt-1 text-xs text-muted-foreground">
                  {formatMoney(state.quote.fine.dailyRateMinor)} per started overdue day. Fine becomes unpaid when returned.
                </p>
              </Alert>
            )}
          </div>
        )}

        <DialogFooter>
          <DialogClose asChild>
            <Button disabled={state.status === "submitting"} variant="outline">Cancel</Button>
          </DialogClose>
          {state.status === "ready" && (
            <Button onClick={() => void returnBook(state.quote)}>Confirm return</Button>
          )}
          {state.status === "submitting" && (
            <Button disabled>
              <LoaderCircle aria-hidden="true" className="animate-spin motion-reduce:animate-none" /> Returning…
            </Button>
          )}
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
