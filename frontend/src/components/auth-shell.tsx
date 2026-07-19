import { ArrowLeft, BookCheck, Clock3, Search } from "lucide-react";
import Link from "next/link";

import { BrandMark } from "@/components/brand-mark";

export function AuthShell({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <div className="min-h-screen lg:grid lg:grid-cols-[minmax(23rem,0.9fr)_minmax(34rem,1.1fr)]">
      <aside className="relative hidden min-h-screen overflow-hidden bg-foreground p-10 text-background lg:flex lg:flex-col lg:justify-between xl:p-14">
        <Link
          href="/"
          className="relative z-10 inline-flex w-fit items-center gap-2 rounded-sm font-display text-2xl font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-background"
        >
          <BrandMark variant="on-dark" />
        </Link>

        <div
          aria-hidden="true"
          className="absolute -right-16 top-1/2 w-72 -translate-y-1/2 rotate-[-4deg]"
        >
          <div className="ml-14 h-12 border border-background/20 bg-primary" />
          <div className="mt-2 h-16 border border-background/20 bg-book-rust" />
          <div className="ml-8 mt-2 h-14 border border-background/20 bg-book-gold" />
          <div className="mt-3 h-1 bg-background/60" />
        </div>

        <div className="relative z-10 max-w-md">
          <p className="font-display text-4xl font-semibold leading-[1.05] xl:text-5xl">
            More reading.
            <span className="block italic text-book-gold">Less keeping track.</span>
          </p>
          <p className="mt-5 max-w-sm text-sm leading-6 text-background/70">
            Your quiet place to find books, borrow online, and follow every loan.
          </p>
          <ul className="mt-8 grid gap-3 border-t border-background/20 pt-5 text-xs font-semibold text-background/80">
            <li className="flex items-center gap-3">
              <Search aria-hidden="true" className="size-4 text-book-gold" />
              Browse available titles
            </li>
            <li className="flex items-center gap-3">
              <BookCheck aria-hidden="true" className="size-4 text-book-gold" />
              Borrow from your account
            </li>
            <li className="flex items-center gap-3">
              <Clock3 aria-hidden="true" className="size-4 text-book-gold" />
              Review loans and history
            </li>
          </ul>
        </div>
      </aside>

      <section className="flex min-h-screen flex-col">
        <header className="border-b border-border px-5 md:px-8">
          <div className="mx-auto flex h-16 max-w-xl items-center justify-between">
            <Link
              href="/"
              className="inline-flex items-center gap-2 rounded-sm font-display text-2xl font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring lg:hidden"
            >
              <BrandMark variant="on-light" />
            </Link>
            <Link
              href="/"
              className="ml-auto inline-flex items-center gap-1 rounded-sm text-xs font-semibold text-muted-foreground outline-none hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring"
            >
              <ArrowLeft aria-hidden="true" className="size-4" />
              Home
            </Link>
          </div>
        </header>
        <main className="flex flex-1 items-center px-5 py-10 sm:px-8 sm:py-14">
          <div className="mx-auto w-full max-w-xl">{children}</div>
        </main>
        <footer className="border-t border-border bg-secondary px-5 py-5 text-center text-xs text-muted-foreground">
          Perpus Digital member access · Private by design
        </footer>
      </section>
    </div>
  );
}
