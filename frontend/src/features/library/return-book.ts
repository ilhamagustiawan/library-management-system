import { z } from "zod";

const fineSchema = z.object({
  overdueDays: z.number().int().positive(),
  dailyRateMinor: z.number().int().positive(),
  totalAmountMinor: z.number().int().positive(),
  currency: z.literal("IDR"),
});
const quoteSchema = z.object({
  status: z.literal("ready"),
  quote: z.object({
    loanId: z.uuid(),
    bookId: z.uuid(),
    dueAt: z.iso.datetime(),
    quotedAt: z.iso.datetime(),
    fine: fineSchema.nullable(),
  }),
});
const resultSchema = z.object({
  status: z.literal("returned"),
  stockUpdate: z.enum(["pending", "confirmed"]),
  fine: z
    .object({
      overdueDays: z.number().int().positive(),
      totalAmountMinor: z.number().int().positive(),
      currency: z.literal("IDR"),
    })
    .nullable(),
});
const errorSchema = z.object({ error: z.object({ message: z.string().min(1) }) });

export type ReturnQuote = z.infer<typeof quoteSchema>["quote"];
export type ReturnResult = z.infer<typeof resultSchema>;

type LoadQuoteResult =
  | { status: "success"; quote: ReturnQuote }
  | { status: "error"; message: string };

type SubmitResult =
  | { status: "success"; result: ReturnResult }
  | { status: "quote-changed"; message: string }
  | { status: "error"; message: string };

async function body(response: Response): Promise<unknown | undefined> {
  try {
    const value: unknown = await response.json();
    return value;
  } catch {
    return undefined;
  }
}

function message(value: unknown, fallback: string) {
  const failure = errorSchema.safeParse(value);
  return failure.success ? failure.data.error.message : fallback;
}

async function loadQuote(loanId: string): Promise<LoadQuoteResult> {
  try {
    const response = await fetch(`/api/member/loans/${encodeURIComponent(loanId)}/return`, {
      headers: { Accept: "application/json" },
      cache: "no-store",
    });
    const value = await body(response);
    const result = quoteSchema.safeParse(value);
    return response.ok && result.success
      ? { status: "success", quote: result.data.quote }
      : { status: "error", message: message(value, "Return quote could not be loaded. Try again.") };
  } catch {
    return { status: "error", message: "Return quote could not be loaded. Check your connection and try again." };
  }
}

async function submit(loanId: string, quote: ReturnQuote): Promise<SubmitResult> {
  try {
    const response = await fetch(`/api/member/loans/${encodeURIComponent(loanId)}/return`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ acceptedFineAmountMinor: quote.fine?.totalAmountMinor ?? 0 }),
    });
    const value = await body(response);
    if (response.status === 409) {
      return {
        status: "quote-changed",
        message: message(value, "Fine changed before return. Review the updated amount and confirm again."),
      };
    }
    const result = resultSchema.safeParse(value);
    return response.ok && result.success
      ? { status: "success", result: result.data }
      : { status: "error", message: message(value, "Book could not be returned. Try again.") };
  } catch {
    return { status: "error", message: "Book could not be returned. Check your connection and try again." };
  }
}

export const ReturnBook = { loadQuote, submit } as const;
