import { ArrowRight } from "lucide-react";
import Link from "next/link";

import { BrandMark } from "@/components/brand-mark";

const footerLinks = [
  { href: "/browse", label: "Browse books" },
  { href: "/dashboard#my-books", label: "My books" },
  { href: "/dashboard#history", label: "Loan history" },
] as const;

export function SiteFooter() {
  return (
    <footer className="border-t-4 border-book-rust bg-foreground text-background">
      <div className="mx-auto grid max-w-6xl gap-10 px-5 py-12 md:grid-cols-[1.35fr_0.65fr_0.8fr] md:px-8 md:py-14">
        <div>
          <Link
            href="/"
            className="inline-flex items-center gap-2 rounded-sm font-display text-3xl font-semibold outline-none focus-visible:ring-2 focus-visible:ring-book-gold"
          >
            <BrandMark variant="on-dark" />
          </Link>
          <p className="mt-4 max-w-sm text-sm leading-6 text-background/65">
            Find available books, borrow online, and keep every loan in one quiet place.
          </p>
        </div>

        <nav aria-label="Reading desk">
          <p className="text-xs font-semibold uppercase tracking-[0.16em] text-book-gold">
            Reading desk
          </p>
          <ul className="mt-4 grid gap-3 text-sm">
            {footerLinks.map((link) => (
              <li key={link.href}>
                <Link
                  href={link.href}
                  className="rounded-sm text-background/70 underline-offset-4 hover:text-background hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-book-gold"
                >
                  {link.label}
                </Link>
              </li>
            ))}
          </ul>
        </nav>

        <div className="border-t border-background/20 pt-6 md:border-l md:border-t-0 md:pl-8 md:pt-0">
          <p className="font-display text-2xl font-semibold leading-tight">Your next read awaits.</p>
          <p className="mt-3 text-sm leading-6 text-background/65">
            Open your member account and return to the shelves.
          </p>
          <Link
            href="/login"
            className="mt-5 inline-flex items-center gap-2 rounded-sm text-sm font-semibold text-book-gold underline-offset-4 hover:underline focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-book-gold"
          >
            Log in to Perpus Digital
            <ArrowRight aria-hidden="true" className="size-4" />
          </Link>
        </div>
      </div>

      <div className="border-t border-background/15">
        <div className="mx-auto flex max-w-6xl flex-col gap-2 px-5 py-5 text-xs text-background/50 sm:flex-row sm:items-center sm:justify-between md:px-8">
          <p>Perpus Digital member portal</p>
          <p>Private by design. Ready after hours.</p>
        </div>
      </div>
    </footer>
  );
}
