import { z } from "zod";

const fineSchema = z.object({
  overdueDays: z.number().int().positive(),
  totalAmountMinor: z.number().int().positive(),
  currency: z.literal("IDR"),
  status: z.literal("unpaid"),
});

const transactionSchema = z.object({
  id: z.string().min(1),
  loanId: z.string().min(1),
  bookId: z.string().min(1),
  type: z.enum(["borrow", "return"]),
  occurredAt: z.iso.datetime(),
  fine: fineSchema.optional(),
});

const transactionPageSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    items: z.array(transactionSchema),
    page: z.number().int().positive(),
    pageSize: z.number().int().positive(),
    totalItems: z.number().int().nonnegative(),
    totalPages: z.number().int().nonnegative(),
  }),
});

const bookSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    id: z.string().min(1),
    title: z.string().min(1),
    author: z.string().min(1),
    coverUrl: z
      .url()
      .refine((value) => ["http:", "https:"].includes(new URL(value).protocol))
      .nullable(),
    publicationYear: z.number().int().nullable(),
  }),
});

type Transaction = z.infer<typeof transactionSchema>;

export type BookReference =
  | {
      status: "available";
      id: string;
      title: string;
      author: string;
      coverUrl: string | null;
      publicationYear: number | null;
    }
  | { status: "unavailable"; id: string };

export type LoanState =
  | { status: "active" }
  | { status: "returned-on-time"; returnedAt: string }
  | {
      status: "returned-late";
      returnedAt: string;
      overdueDays: number;
      fine: { amountMinor: number; currency: "IDR"; status: "unpaid" };
    };

export type MemberLoan = {
  loanId: string;
  book: BookReference;
  borrowedAt: string;
  state: LoanState;
};

export type MemberLibrary = {
  activeLoans: MemberLoan[];
  history: MemberLoan[];
  summary: {
    activeLoans: number;
    completedLoans: number;
    lateReturns: number;
    unpaidFineMinor: number;
    fineCurrency: "IDR";
  };
};

export type LoadMemberLibraryResult =
  | { status: "success"; library: MemberLibrary }
  | { status: "error"; error: { kind: "unauthorized" | "unavailable" | "invalid-response" } };

type LoadInput = {
  issuer: string;
  accessToken: string;
  fetcher?: typeof fetch;
};

function request(url: URL, accessToken: string, fetcher: typeof fetch) {
  return fetcher(url, {
    headers: { Accept: "application/json", Authorization: `Bearer ${accessToken}` },
    cache: "no-store",
  });
}

async function loadTransactions(input: LoadInput, fetcher: typeof fetch) {
  const firstURL = new URL("/api/v1/transactions/me", input.issuer);
  firstURL.searchParams.set("page", "1");
  firstURL.searchParams.set("pageSize", "100");

  let firstResponse: Response;
  try {
    firstResponse = await request(firstURL, input.accessToken, fetcher);
  } catch {
    return { status: "error", error: { kind: "unavailable" } } as const;
  }
  if (firstResponse.status === 401) {
    return { status: "error", error: { kind: "unauthorized" } } as const;
  }
  if (!firstResponse.ok) {
    return { status: "error", error: { kind: "unavailable" } } as const;
  }

  try {
    const first = transactionPageSchema.safeParse(await firstResponse.json());
    if (!first.success) {
      return { status: "error", error: { kind: "invalid-response" } } as const;
    }
    const remaining = await Promise.all(
      Array.from({ length: Math.max(0, first.data.data.totalPages - 1) }, async (_, index) => {
        const url = new URL(firstURL);
        url.searchParams.set("page", String(index + 2));
        const response = await request(url, input.accessToken, fetcher);
        if (!response.ok) return undefined;
        const parsed = transactionPageSchema.safeParse(await response.json());
        return parsed.success ? parsed.data.data.items : undefined;
      }),
    );
    if (remaining.some((page) => page === undefined)) {
      return { status: "error", error: { kind: "invalid-response" } } as const;
    }
    return {
      status: "success",
      transactions: [
        ...first.data.data.items,
        ...remaining.flatMap((page) => page ?? []),
      ],
    } as const;
  } catch {
    return { status: "error", error: { kind: "invalid-response" } } as const;
  }
}

function buildLoans(transactions: readonly Transaction[]) {
  const grouped = new Map<string, { borrow?: Transaction; returned?: Transaction }>();
  for (const transaction of transactions) {
    const loan = grouped.get(transaction.loanId) ?? {};
    if (transaction.type === "borrow") loan.borrow = transaction;
    else loan.returned = transaction;
    grouped.set(transaction.loanId, loan);
  }

  const loans: Array<Omit<MemberLoan, "book"> & { bookId: string }> = [];
  for (const [loanId, transaction] of grouped) {
    if (transaction.borrow === undefined) return { status: "invalid" } as const;
    const returned = transaction.returned;
    let state: LoanState = { status: "active" };
    if (returned !== undefined) {
      state =
        returned.fine === undefined
          ? { status: "returned-on-time", returnedAt: returned.occurredAt }
          : {
              status: "returned-late",
              returnedAt: returned.occurredAt,
              overdueDays: returned.fine.overdueDays,
              fine: {
                amountMinor: returned.fine.totalAmountMinor,
                currency: returned.fine.currency,
                status: returned.fine.status,
              },
            };
    }
    loans.push({
      loanId,
      bookId: transaction.borrow.bookId,
      borrowedAt: transaction.borrow.occurredAt,
      state,
    });
  }
  loans.sort((left, right) => {
    const leftDate = left.state.status === "active" ? left.borrowedAt : left.state.returnedAt;
    const rightDate = right.state.status === "active" ? right.borrowedAt : right.state.returnedAt;
    return rightDate.localeCompare(leftDate);
  });
  return { status: "success", loans } as const;
}

async function loadBook(
  issuer: string,
  accessToken: string,
  id: string,
  fetcher: typeof fetch,
): Promise<BookReference> {
  try {
    const response = await request(new URL(`/api/v1/books/${encodeURIComponent(id)}`, issuer), accessToken, fetcher);
    if (!response.ok) return { status: "unavailable", id };
    const parsed = bookSchema.safeParse(await response.json());
    return parsed.success
      ? { status: "available", ...parsed.data.data }
      : { status: "unavailable", id };
  } catch {
    return { status: "unavailable", id };
  }
}

async function load(input: LoadInput): Promise<LoadMemberLibraryResult> {
  const fetcher = input.fetcher ?? fetch;
  const transactionResult = await loadTransactions(input, fetcher);
  if (transactionResult.status === "error") return transactionResult;
  const loanResult = buildLoans(transactionResult.transactions);
  if (loanResult.status === "invalid") {
    return { status: "error", error: { kind: "invalid-response" } };
  }

  // TODO: Use a current-loans endpoint once it exposes dueAt; transaction history has no due dates.
  const active = loanResult.loans.filter((loan) => loan.state.status === "active");
  const completed = loanResult.loans.filter((loan) => loan.state.status !== "active");
  const history = [...active, ...completed.slice(0, Math.max(0, 20 - active.length))].sort(
    (left, right) => {
      const leftDate = left.state.status === "active" ? left.borrowedAt : left.state.returnedAt;
      const rightDate = right.state.status === "active" ? right.borrowedAt : right.state.returnedAt;
      return rightDate.localeCompare(leftDate);
    },
  );
  const bookIDs = [...new Set([...active, ...history].map((loan) => loan.bookId))];
  const books = new Map(
    await Promise.all(
      bookIDs.map(async (id) => [id, await loadBook(input.issuer, input.accessToken, id, fetcher)] as const),
    ),
  );
  const withBooks = (loan: (typeof loanResult.loans)[number]): MemberLoan => ({
    loanId: loan.loanId,
    borrowedAt: loan.borrowedAt,
    state: loan.state,
    book: books.get(loan.bookId) ?? { status: "unavailable", id: loan.bookId },
  });
  const late = completed.filter((loan) => loan.state.status === "returned-late");

  return {
    status: "success",
    library: {
      activeLoans: active.map(withBooks),
      history: history.map(withBooks),
      summary: {
        activeLoans: active.length,
        completedLoans: completed.length,
        lateReturns: late.length,
        unpaidFineMinor: late.reduce(
          (total, loan) => total + (loan.state.status === "returned-late" ? loan.state.fine.amountMinor : 0),
          0,
        ),
        fineCurrency: "IDR",
      },
    },
  };
}

function activeLoanForBook(result: LoadMemberLibraryResult, bookID: string) {
  if (result.status === "error") return undefined;
  return result.library.activeLoans.find((loan) => loan.book.id === bookID);
}

export const MemberLibrary = { activeLoanForBook, load } as const;
