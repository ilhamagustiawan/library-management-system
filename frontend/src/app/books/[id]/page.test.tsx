import { render, screen } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { WebSession } from "@/features/auth/web-session";

import BookDetailPage from "./page";

const next = vi.hoisted(() => ({
  activeLoanForBook: vi.fn(),
  getBook: vi.fn(),
  loadLibrary: vi.fn(),
  refresh: vi.fn(),
}));

vi.mock("next/headers", () => ({
  cookies: async () => ({ get: () => ({ value: nextCookie }) }),
}));
vi.mock("next/navigation", () => ({
  notFound: vi.fn((): never => { throw new Error("not-found"); }),
  redirect: vi.fn((target: string): never => { throw new Error(`redirect:${target}`); }),
  usePathname: () => "/books/7b36fe43-f31d-4861-884f-42ed7386b1e9",
  useRouter: () => ({ refresh: next.refresh }),
}));
vi.mock("@/features/library/book-catalog", () => ({ BookCatalog: { get: next.getBook } }));
vi.mock("@/features/library/member-library", () => ({
  MemberLibrary: { activeLoanForBook: next.activeLoanForBook, load: next.loadLibrary },
}));

const sessionSecret = "abcdef0123456789abcdef0123456789";
let nextCookie = "";

describe("BookDetailPage", () => {
  beforeEach(() => {
    vi.stubEnv("AUTH_ISSUER", "http://localhost:8000");
    vi.stubEnv("AUTH_CLIENT_ID", "member-nextjs-web");
    vi.stubEnv("AUTH_CLIENT_SECRET", "0123456789abcdef0123456789abcdef");
    vi.stubEnv("AUTH_REDIRECT_URI", "http://localhost:3000/api/auth/callback/library");
    vi.stubEnv("AUTH_SESSION_SECRET", sessionSecret);
    nextCookie = WebSession.seal(
      {
        user: { id: "member-1", name: "Maya Chen", email: "maya@perpus-digital.test" },
        accessToken: "access-token",
        refreshToken: "refresh-token",
        tokenType: "Bearer",
        scope: "books:read loans:return:self",
        expiresAt: Math.floor(Date.now() / 1_000) + 900,
      },
      sessionSecret,
    );
    next.getBook.mockResolvedValue({
      status: "success",
      book: {
        id: "7b36fe43-f31d-4861-884f-42ed7386b1e9",
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
      },
    });
    next.loadLibrary.mockResolvedValue({ status: "success", library: {} });
    next.activeLoanForBook.mockReturnValue({
      loanId: "52a88672-a4c2-4876-be5a-65863aeb35e4",
    });
  });

  it("offers return instead of borrow for the member's active loan", async () => {
    render(
      await BookDetailPage({
        params: Promise.resolve({ id: "7b36fe43-f31d-4861-884f-42ed7386b1e9" }),
      }),
    );

    expect(screen.getByRole("button", { name: "Return this book" })).toBeVisible();
    expect(screen.queryByRole("button", { name: "Borrow this book" })).not.toBeInTheDocument();
  });
});
