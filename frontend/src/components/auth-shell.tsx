import { ArrowLeft, BookOpen, Quote } from "lucide-react";
import Link from "next/link";

export function AuthShell({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <main className="grid min-h-screen lg:grid-cols-[0.85fr_1.15fr]">
      <section className="relative hidden overflow-hidden bg-primary p-12 text-primary-foreground lg:flex lg:flex-col lg:justify-between">
        <Link
          href="/"
          className="inline-flex w-fit items-center gap-2 rounded-sm font-display text-2xl font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-foreground"
        >
          <span className="grid size-8 place-items-center rounded-full bg-primary-foreground text-primary">
            <BookOpen aria-hidden="true" className="size-4" />
          </span>
          Libry
        </Link>
        <div className="absolute -right-28 top-1/2 size-96 -translate-y-1/2 rounded-full border border-primary-foreground/20" />
        <div className="absolute -right-10 top-1/2 size-56 -translate-y-1/2 rounded-full border border-primary-foreground/20" />
        <figure className="relative max-w-md">
          <Quote aria-hidden="true" className="mb-6 size-10 opacity-60" strokeWidth={1.4} />
          <blockquote className="font-display text-4xl font-semibold leading-tight">
            A library card is a small key to a very large world.
          </blockquote>
          <figcaption className="mt-5 text-sm text-primary-foreground/70">
            Your member space keeps that world within reach.
          </figcaption>
        </figure>
      </section>
      <section className="flex min-h-screen flex-col px-5 py-6 sm:px-8 lg:px-14">
        <div className="flex items-center justify-between">
          <Link
            href="/"
            className="inline-flex items-center gap-2 rounded-sm text-sm font-semibold text-muted-foreground outline-none hover:text-foreground focus-visible:ring-2 focus-visible:ring-ring lg:hidden"
          >
            <ArrowLeft aria-hidden="true" className="size-4" />
            Home
          </Link>
          <Link href="/" className="ml-auto hidden font-display text-2xl font-semibold lg:block">
            Libry
          </Link>
        </div>
        <div className="my-auto mx-auto w-full max-w-xl py-10">{children}</div>
      </section>
    </main>
  );
}
