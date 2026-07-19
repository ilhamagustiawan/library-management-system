"use client";

import { ArrowRight, Search, SearchX, X } from "lucide-react";
import Link from "next/link";
import { useEffect, useMemo, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Empty,
  EmptyContent,
  EmptyDescription,
  EmptyHeader,
  EmptyMedia,
  EmptyTitle,
} from "@/components/ui/empty";
import {
  InputGroup,
  InputGroupAddon,
  InputGroupButton,
  InputGroupInput,
} from "@/components/ui/input-group";

import type { CatalogBook } from "./book-catalog";
import { BookCover } from "./book-cover";

export function BookBrowser({ books }: { books: CatalogBook[] }) {
  const [query, setQuery] = useState("");
  const [entranceComplete, setEntranceComplete] = useState(false);
  const searchRef = useRef<HTMLInputElement>(null);
  useEffect(() => {
    const timer = window.setTimeout(() => setEntranceComplete(true), 500);
    return () => window.clearTimeout(timer);
  }, []);

  const visibleBooks = useMemo(() => {
    const normalized = query.trim().toLocaleLowerCase();
    if (normalized === "") return books;
    return books.filter((book) =>
      `${book.title} ${book.author} ${book.isbn}`.toLocaleLowerCase().includes(normalized),
    );
  }, [books, query]);

  function clearSearch() {
    setQuery("");
    searchRef.current?.focus();
  }

  return (
    <>
      <div className="mt-7 max-w-2xl">
        <label htmlFor="book-search" className="sr-only">Search available books</label>
        <InputGroup>
          <InputGroupInput
            id="book-search"
            ref={searchRef}
            type="search"
            value={query}
            onChange={(event) => setQuery(event.target.value)}
            placeholder="Search title, author, or ISBN"
            className="text-base [&::-webkit-search-cancel-button]:appearance-none"
          />
          <InputGroupAddon>
            <Search aria-hidden="true" />
          </InputGroupAddon>
          {query !== "" && (
            <InputGroupAddon align="inline-end">
              <InputGroupButton size="icon-xs" aria-label="Clear search" onClick={clearSearch}>
                <X aria-hidden="true" />
              </InputGroupButton>
            </InputGroupAddon>
          )}
        </InputGroup>
      </div>
      <p aria-live="polite" className="mt-4 text-sm text-muted-foreground">
        {visibleBooks.length} {visibleBooks.length === 1 ? "book" : "books"} available
      </p>
      {visibleBooks.length === 0 ? (
        <Empty className="mt-8">
          <EmptyHeader>
            <EmptyMedia variant="icon">
              <SearchX aria-hidden="true" />
            </EmptyMedia>
            <EmptyTitle>No matching books</EmptyTitle>
            <EmptyDescription>Try another title, author, or ISBN.</EmptyDescription>
          </EmptyHeader>
          <EmptyContent>
            <Button variant="outline" onClick={clearSearch}>
              Clear search
            </Button>
          </EmptyContent>
        </Empty>
      ) : (
        <ul className="mt-5 grid grid-cols-2 gap-x-5 gap-y-9 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5">
          {visibleBooks.map((book, index) => {
            const entranceIndex = Math.min(index, 7);
            return (
              <li
                key={book.id}
                data-enter={entranceComplete ? undefined : String(entranceIndex)}
                className="catalog-card group min-w-0"
              >
                <Link
                  href={`/books/${book.id}`}
                  className="book-card block rounded-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary focus-visible:ring-offset-4"
                >
                  <div className="book-card-cover">
                    <BookCover coverUrl={book.coverUrl} title={book.title} />
                  </div>
                  <h2 className="mt-4 line-clamp-2 font-display text-lg font-semibold leading-tight decoration-primary underline-offset-4 group-hover:underline sm:text-xl">
                    {book.title}
                  </h2>
                  <p className="mt-1 truncate text-sm text-muted-foreground">{book.author}</p>
                  <p className="mt-2 flex items-center gap-1 text-xs font-bold uppercase tracking-[0.12em] text-primary">
                    {book.availableCopies} available
                    <ArrowRight aria-hidden="true" className="book-card-arrow size-3.5" />
                  </p>
                </Link>
              </li>
            );
          })}
        </ul>
      )}
    </>
  );
}
