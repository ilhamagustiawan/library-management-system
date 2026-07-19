import { ArrowRight, Bookmark, Clock3, LibraryBig } from "lucide-react";
import Link from "next/link";

import { SiteFooter } from "@/components/site-footer";
import { SiteHeader } from "@/components/site-header";
import { buttonVariants } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const benefits = [
  {
    number: "01",
    icon: Bookmark,
    title: "Keep your place",
    copy: "A clear home for the books, holds, and reading plans tied to your membership.",
  },
  {
    number: "02",
    icon: Clock3,
    title: "Know what is next",
    copy: "See important borrowing details without searching through receipts or emails.",
  },
  {
    number: "03",
    icon: LibraryBig,
    title: "Stay close to the library",
    copy: "One calm doorway to the member services your library makes available.",
  },
] as const;

export default function HomePage() {
  return (
    <div className="min-h-screen">
      <a
        href="#main-content"
        className="sr-only z-50 bg-primary px-4 py-2 text-primary-foreground focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
      >
        Skip to content
      </a>
      <SiteHeader />
      <main id="main-content">
        <section className="mx-auto grid max-w-7xl gap-14 px-5 py-16 md:px-8 md:py-24 lg:grid-cols-[1.15fr_0.85fr] lg:items-center lg:py-32">
          <div>
            <p className="mb-6 flex items-center gap-3 text-xs font-semibold uppercase tracking-[0.24em] text-primary">
              <span className="h-px w-10 bg-primary" />
              Your member desk, online
            </p>
            <h1 className="max-w-3xl font-display text-6xl font-semibold leading-[0.94] tracking-[-0.04em] sm:text-7xl lg:text-8xl">
              Your library,
              <span className="block italic text-primary">after hours.</span>
            </h1>
            <p className="mt-7 max-w-xl text-lg leading-8 text-muted-foreground">
              Libry gives members one thoughtful place to stay connected with their library—simple,
              quiet, and ready when they are.
            </p>
            <div className="mt-9 flex flex-col gap-3 sm:flex-row">
              <Link href="/register" className={cn(buttonVariants({ size: "lg" }))}>
                Create member account
                <ArrowRight aria-hidden="true" className="size-4" />
              </Link>
              <Link
                href="/login"
                className={cn(buttonVariants({ variant: "outline", size: "lg" }))}
              >
                I already have an account
              </Link>
            </div>
          </div>

          <div aria-hidden="true" className="relative mx-auto w-full max-w-md lg:mr-0">
            <div className="absolute -left-7 top-10 hidden h-72 w-72 rounded-full border border-primary/25 sm:block" />
            <div className="relative ml-auto aspect-[4/5] w-[88%] overflow-hidden rounded-t-[10rem] border border-border bg-secondary px-8 pb-9 pt-20 shadow-[10px_12px_0_0_var(--primary)]">
              <div className="absolute inset-x-8 top-9 flex items-center justify-between border-b border-border pb-3 text-xs font-semibold uppercase tracking-[0.18em] text-muted-foreground">
                <span>Member no. 0248</span>
                <span>Libry</span>
              </div>
              <div className="flex h-full items-end justify-center gap-2 border-b-4 border-foreground/75">
                <div className="h-[58%] w-14 rounded-t bg-primary" />
                <div className="grid h-[78%] w-16 place-items-center rounded-t border-2 border-foreground/70 bg-accent">
                  <span className="-rotate-90 whitespace-nowrap font-display text-sm font-semibold uppercase tracking-widest">
                    Read often
                  </span>
                </div>
                <div className="h-[68%] w-12 rounded-t bg-foreground/80" />
                <div className="h-[88%] w-16 rounded-t border-2 border-primary bg-background" />
              </div>
            </div>
          </div>
        </section>

        <section className="border-y border-border/70 bg-card/65">
          <div className="mx-auto max-w-7xl px-5 py-16 md:px-8 md:py-24">
            <div className="grid gap-8 border-b border-border pb-10 lg:grid-cols-2">
              <h2 className="font-display text-4xl font-semibold tracking-tight sm:text-5xl">
                Membership should feel
                <span className="italic text-primary"> effortless.</span>
              </h2>
              <p className="max-w-xl text-base leading-7 text-muted-foreground lg:justify-self-end">
                Built for readers who want less administration between them and their next visit.
              </p>
            </div>
            <ol>
              {benefits.map(({ copy, icon: Icon, number, title }) => (
                <li
                  key={number}
                  className="grid gap-4 border-b border-border py-8 sm:grid-cols-[4rem_3rem_1fr] sm:items-start lg:grid-cols-[6rem_4rem_1fr_1fr]"
                >
                  <span className="text-sm font-semibold text-primary">{number}</span>
                  <Icon aria-hidden="true" className="size-6 text-primary" strokeWidth={1.6} />
                  <h3 className="font-display text-2xl font-semibold">{title}</h3>
                  <p className="max-w-lg leading-7 text-muted-foreground lg:justify-self-end">{copy}</p>
                </li>
              ))}
            </ol>
          </div>
        </section>
      </main>
      <SiteFooter />
    </div>
  );
}
