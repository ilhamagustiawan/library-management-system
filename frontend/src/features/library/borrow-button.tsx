"use client";

import { ArrowRight, LoaderCircle } from "lucide-react";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { z } from "zod";

import { Button } from "@/components/ui/button";

import type { LoanEligibility } from "./loan-eligibility";

const errorSchema = z.object({
  error: z.object({
    kind: z.enum([
      "invalid-input",
      "session-expired",
      "loan-limit",
      "already-borrowed",
      "unavailable",
      "not-found",
      "service-unavailable",
    ]),
    message: z.string().min(1),
  }),
});

type State =
  | { status: "idle" }
  | { status: "pending" }
  | { status: "error"; message: string };

export function BorrowButton({ bookId, eligibility }: { bookId: string; eligibility: LoanEligibility }) {
  const router = useRouter();
  const [state, setState] = useState<State>({ status: "idle" });
  const disabled = eligibility.status === "disabled" || state.status === "pending";

  async function borrow() {
    setState({ status: "pending" });
    try {
      const response = await fetch("/api/member/loans", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ bookId }),
      });
      if (response.ok) {
        router.push("/dashboard?borrowed=1");
        return;
      }
      const result = errorSchema.safeParse(await response.json());
      setState({
        status: "error",
        message: result.success ? result.data.error.message : "Book could not be borrowed. Try again.",
      });
    } catch {
      setState({ status: "error", message: "Book could not be borrowed. Check your connection and try again." });
    }
  }

  return (
    <div>
      <Button className="h-12 w-full text-base" disabled={disabled} onClick={borrow}>
        {state.status === "pending" ? (
          <>
            <LoaderCircle
              aria-hidden="true"
              data-icon="inline-start"
              className="animate-spin motion-reduce:animate-none"
            />
            Borrowing…
          </>
        ) : (
          <>
            Borrow this book <ArrowRight aria-hidden="true" data-icon="inline-end" />
          </>
        )}
      </Button>
      {eligibility.status === "disabled" && (
        <p className="mt-3 text-sm leading-5 text-muted-foreground">{eligibility.message}</p>
      )}
      {state.status === "error" && (
        <p role="alert" className="mt-3 text-sm leading-5 text-destructive">{state.message}</p>
      )}
    </div>
  );
}
