import { z } from "zod";

const bookSchema = z.object({
  id: z.uuid(),
  isbn: z.string().min(1),
  title: z.string().min(1),
  author: z.string().min(1),
  description: z.string().nullable(),
  coverUrl: z
    .url()
    .refine((value) => ["http:", "https:"].includes(new URL(value).protocol))
    .nullable(),
  publicationYear: z.number().int().nullable(),
  totalCopies: z.number().int().nonnegative(),
  availableCopies: z.number().int().nonnegative(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

const listSchema = z.object({
  code: z.literal("LMS-200000"),
  data: z.object({
    items: z.array(bookSchema),
    pagination: z.object({
      page: z.number().int().positive(),
      pageSize: z.number().int().positive(),
      totalItems: z.number().int().nonnegative(),
      totalPages: z.number().int().nonnegative(),
    }),
  }),
});

const detailSchema = z.object({ code: z.literal("LMS-200000"), data: bookSchema });

export type CatalogBook = z.infer<typeof bookSchema>;

export type CatalogError = {
  kind: "unauthorized" | "not-found" | "unavailable" | "invalid-response";
};

export type LoadCatalogResult =
  | { status: "success"; books: CatalogBook[] }
  | { status: "error"; error: CatalogError };

export type LoadBookResult =
  | { status: "success"; book: CatalogBook }
  | { status: "error"; error: CatalogError };

type CatalogInput = {
  issuer: string;
  accessToken: string;
  fetcher?: typeof fetch;
};

function request(url: URL, input: CatalogInput, fetcher: typeof fetch) {
  return fetcher(url, {
    headers: { Accept: "application/json", Authorization: `Bearer ${input.accessToken}` },
    cache: "no-store",
  });
}

function responseError(status: number): CatalogError {
  if (status === 401) return { kind: "unauthorized" };
  if (status === 404) return { kind: "not-found" };
  return { kind: "unavailable" };
}

async function listAvailable(input: CatalogInput): Promise<LoadCatalogResult> {
  const fetcher = input.fetcher ?? fetch;
  const firstUrl = new URL("/api/v1/books", input.issuer);
  firstUrl.searchParams.set("availableOnly", "true");
  firstUrl.searchParams.set("page", "1");
  firstUrl.searchParams.set("pageSize", "100");
  firstUrl.searchParams.set("sortBy", "title");
  firstUrl.searchParams.set("sortOrder", "asc");

  try {
    const firstResponse = await request(firstUrl, input, fetcher);
    if (!firstResponse.ok) return { status: "error", error: responseError(firstResponse.status) };
    const first = listSchema.safeParse(await firstResponse.json());
    if (!first.success) return { status: "error", error: { kind: "invalid-response" } };

    const remaining = await Promise.all(
      Array.from({ length: Math.max(0, first.data.data.pagination.totalPages - 1) }, async (_, index) => {
        const url = new URL(firstUrl);
        url.searchParams.set("page", String(index + 2));
        const response = await request(url, input, fetcher);
        if (!response.ok) return { status: "error", error: responseError(response.status) } as const;
        const page = listSchema.safeParse(await response.json());
        return page.success
          ? ({ status: "success", books: page.data.data.items } as const)
          : ({ status: "error", error: { kind: "invalid-response" } } as const);
      }),
    );
    const failedPage = remaining.find((page) => page.status === "error");
    if (failedPage !== undefined) return failedPage;

    return {
      status: "success",
      books: [
        ...first.data.data.items,
        ...remaining.flatMap((page) => (page.status === "success" ? page.books : [])),
      ],
    };
  } catch {
    return { status: "error", error: { kind: "invalid-response" } };
  }
}

async function get(input: CatalogInput & { id: string }): Promise<LoadBookResult> {
  const id = z.uuid().safeParse(input.id);
  if (!id.success) return { status: "error", error: { kind: "not-found" } };
  const fetcher = input.fetcher ?? fetch;

  try {
    const response = await request(
      new URL(`/api/v1/books/${encodeURIComponent(id.data)}`, input.issuer),
      input,
      fetcher,
    );
    if (!response.ok) return { status: "error", error: responseError(response.status) };
    const parsed = detailSchema.safeParse(await response.json());
    return parsed.success
      ? { status: "success", book: parsed.data.data }
      : { status: "error", error: { kind: "invalid-response" } };
  } catch {
    return { status: "error", error: { kind: "invalid-response" } };
  }
}

export const BookCatalog = { get, listAvailable } as const;
