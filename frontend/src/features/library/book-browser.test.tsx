import { act, render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";

import type { CatalogBook } from "./book-catalog";
import { BookBrowser } from "./book-browser";

function book(id: string, title: string, author: string): CatalogBook {
  return {
    id,
    isbn: "9780132350884",
    title,
    author,
    description: null,
    coverUrl: null,
    publicationYear: 2008,
    totalCopies: 3,
    availableCopies: 2,
    createdAt: "2026-07-01T00:00:00Z",
    updatedAt: "2026-07-01T00:00:00Z",
  };
}

const books = [
  book("0ec82798-8ff9-48c5-b68f-2b8c050647ac", "Clean Code", "Robert C. Martin"),
  book("343b5a45-6419-4e5b-b34b-724ef74f0db9", "The Pragmatic Programmer", "David Thomas"),
];

describe("BookBrowser", () => {
  afterEach(() => vi.useRealTimers());

  it("filters books and clears an empty search while restoring focus", async () => {
    const user = userEvent.setup();
    render(<BookBrowser books={books} />);

    const search = screen.getByRole("searchbox", { name: "Search available books" });
    await user.type(search, "missing title");

    expect(screen.getByText("No matching books")).toBeVisible();
    const clearActions = screen.getAllByRole("button", { name: "Clear search" });
    const emptyStateAction = clearActions.at(-1);
    if (emptyStateAction === undefined) throw new Error("Expected empty-state clear action.");
    await user.click(emptyStateAction);

    expect(search).toHaveValue("");
    expect(search).toHaveFocus();
    expect(screen.getByText("2 books available")).toBeVisible();
  });

  it("runs the grid entrance once, then leaves filtered results still", async () => {
    vi.useFakeTimers();
    render(<BookBrowser books={books} />);

    const firstCard = screen.getByRole("link", { name: /clean code/i }).closest("li");
    const lastCard = screen.getByRole("link", { name: /pragmatic programmer/i }).closest("li");
    expect(firstCard).toHaveAttribute("data-enter", "0");
    expect(lastCard).toHaveAttribute("data-enter", "1");
    if (lastCard === null) throw new Error("Expected last catalog card.");

    act(() => vi.advanceTimersByTime(500));

    expect(firstCard).not.toHaveAttribute("data-enter");
    expect(lastCard).not.toHaveAttribute("data-enter");
  });
});
