import { describe, expect, it, vi } from "vitest";

import { BookCatalog } from "./book-catalog";

const issuer = "http://localhost:8000";

function book(id: string, title: string) {
  return {
    id,
    isbn: "9780132350884",
    title,
    author: "Robert C. Martin",
    description: null,
    coverUrl: null,
    publicationYear: 2008,
    totalCopies: 3,
    availableCopies: 2,
    createdAt: "2026-07-01T00:00:00Z",
    updatedAt: "2026-07-01T00:00:00Z",
  };
}

describe("BookCatalog", () => {
  it("loads every page of available books", async () => {
    const fetcher = vi.fn<typeof fetch>(async (input) => {
      const url = new URL(input.toString());
      const page = Number(url.searchParams.get("page"));
      return Response.json({
        code: "LMS-200000",
        data: {
          items: page === 1
            ? [book("0ec82798-8ff9-48c5-b68f-2b8c050647ac", "Clean Code")]
            : [book("343b5a45-6419-4e5b-b34b-724ef74f0db9", "The Pragmatic Programmer")],
          pagination: { page, pageSize: 100, totalItems: 2, totalPages: 2 },
        },
      });
    });

    const result = await BookCatalog.listAvailable({ issuer, accessToken: "access-token", fetcher });

    expect(result.status).toBe("success");
    if (result.status !== "success") return;
    expect(result.books.map(({ title }) => title)).toEqual(["Clean Code", "The Pragmatic Programmer"]);
    expect(fetcher).toHaveBeenCalledTimes(2);
    expect(fetcher.mock.calls[0]?.[0].toString()).toContain("availableOnly=true");
  });

  it("rejects invalid detail payloads", async () => {
    const result = await BookCatalog.get({
      issuer,
      accessToken: "access-token",
      id: "0ec82798-8ff9-48c5-b68f-2b8c050647ac",
      fetcher: vi.fn<typeof fetch>(async () => Response.json({ code: "LMS-200000", data: {} })),
    });

    expect(result).toEqual({ status: "error", error: { kind: "invalid-response" } });
  });

  it("rejects cover URLs outside HTTP and HTTPS", async () => {
    const payload = book("0ec82798-8ff9-48c5-b68f-2b8c050647ac", "Clean Code");
    payload.coverUrl = "data:image/png;base64,unsafe";
    const result = await BookCatalog.get({
      issuer,
      accessToken: "access-token",
      id: payload.id,
      fetcher: vi.fn<typeof fetch>(async () =>
        Response.json({ code: "LMS-200000", data: payload }),
      ),
    });

    expect(result).toEqual({ status: "error", error: { kind: "invalid-response" } });
  });
});
