import { ArrowRight, BookCheck, Clock3, Search } from "lucide-react";
import { cookies } from "next/headers";
import Image from "next/image";
import Link from "next/link";

import { SiteFooter } from "@/components/site-footer";
import { SiteHeader, type SiteHeaderAuth } from "@/components/site-header";
import { buttonVariants } from "@/components/ui/button";
import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { WebSession } from "@/features/auth/web-session";
import { cn } from "@/lib/utils";

export const dynamic = "force-dynamic";

const benefits = [
  {
    number: "01",
    icon: Search,
    title: "Know what is available before you visit",
    copy: "Browse books currently available and open any title for its full details.",
  },
  {
    number: "02",
    icon: BookCheck,
    title: "Borrow from wherever you are",
    copy: "Choose an available book and add it to your account in a few clear steps.",
  },
  {
    number: "03",
    icon: Clock3,
    title: "See every loan in one place",
    copy: "Check current books, return status, borrowing history, and any fines together.",
  },
] as const;

async function loadHeaderAuth(): Promise<SiteHeaderAuth> {
  const sealedSession = (await cookies()).get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) return { status: "guest" };

  const config = AuthConfig.load();
  const session = WebSession.open(sealedSession, config.sessionSecret);
  return session.status === "valid"
    ? { status: "authenticated", logoutEndpoint: config.logoutEndpoint }
    : { status: "guest" };
}

export default async function HomePage() {
  const auth = await loadHeaderAuth();

  return (
    <div className="flex min-h-screen flex-col">
      <a
        href="#main-content"
        className="sr-only z-50 bg-primary px-4 py-2 text-primary-foreground focus:not-sr-only focus:fixed focus:left-4 focus:top-4"
      >
        Skip to content
      </a>
      <SiteHeader auth={auth} />
      <main id="main-content" className="flex-1">
        <section className="border-b border-border bg-secondary">
          <div className="mx-auto max-w-6xl px-5 py-14 md:px-8 md:py-20 lg:py-24">
            <div className="landing-copy max-w-4xl">
              <p className="mb-4 text-xs font-bold uppercase tracking-[0.2em] text-primary">
                The library, on your time
              </p>
              <h1 className="max-w-3xl font-display text-5xl font-semibold leading-[0.98] tracking-[-0.035em] sm:text-6xl lg:text-7xl">
                Find your next read.{" "}
                <span className="block italic text-book-rust">Keep every loan in view.</span>
              </h1>
              <p className="mt-6 max-w-xl text-base leading-7 text-muted-foreground sm:text-lg sm:leading-8">
                Browse available books, borrow online, and follow your reading history from one
                calm member space.
              </p>
              <div className="mt-8 flex flex-col gap-3 sm:flex-row">
                {auth.status === "authenticated" ? (
                  <Link
                    href="/dashboard"
                    className={cn(buttonVariants({ size: "lg" }), "sm:min-w-48")}
                  >
                    Open member home
                    <ArrowRight aria-hidden="true" data-icon="inline-end" />
                  </Link>
                ) : (
                  <>
                    <Link
                      href="/register"
                      className={cn(buttonVariants({ size: "lg" }), "sm:min-w-48")}
                    >
                      Create my member account
                      <ArrowRight aria-hidden="true" data-icon="inline-end" />
                    </Link>
                    <Link
                      href="/login"
                      className={cn(buttonVariants({ variant: "outline", size: "lg" }))}
                    >
                      Log in to my account
                    </Link>
                  </>
                )}
              </div>
            </div>

            <figure className="landing-artwork mt-10 w-full md:mt-12">
              <div className="border border-foreground/30 bg-card p-2 shadow-[8px_8px_0_var(--accent)]">
                <div className="relative aspect-[4/3] overflow-hidden bg-muted sm:aspect-video lg:aspect-[21/9]">
                  <Image
                    src="/images/perpus-digital-hero-reading-desk-banner.png"
                    alt="An open book on a library reading desk surrounded by shelves"
                    fill
                    priority
                    sizes="(min-width: 1280px) 72rem, calc(100vw - 2.5rem)"
                    className="object-cover object-center"
                  />
                </div>
                <figcaption className="flex items-center justify-between gap-4 bg-foreground px-4 py-3 text-xs text-background">
                  <span className="font-semibold uppercase tracking-[0.12em]">Your reading desk</span>
                  <span className="text-background/70">Find · Borrow · Follow</span>
                </figcaption>
              </div>
            </figure>
          </div>
        </section>

        <section className="bg-background">
          <div className="mx-auto grid max-w-6xl gap-12 px-5 py-14 md:px-8 md:py-20 lg:grid-cols-[1fr_19rem]">
            <div>
              <div className="border-b border-border pb-5">
                <h2 className="font-display text-4xl font-semibold leading-tight tracking-[-0.025em] sm:text-5xl">
                  Find it. Borrow it. Keep track of it.
                </h2>
                <p className="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
                  Everything needed to choose, borrow, and keep track—without the paper trail.
                </p>
              </div>
              <ol>
                {benefits.map(({ copy, icon: Icon, number, title }) => (
                  <li
                    key={number}
                    className="grid gap-3 border-b border-border py-5 sm:grid-cols-[2.5rem_2rem_1fr] sm:items-start lg:grid-cols-[2.5rem_2rem_12rem_1fr]"
                  >
                    <span className="text-xs font-semibold text-muted-foreground">{number}</span>
                    <Icon aria-hidden="true" className="size-5 text-book-rust" strokeWidth={1.8} />
                    <h3 className="font-display text-xl font-semibold">{title}</h3>
                    <p className="text-sm leading-6 text-muted-foreground">{copy}</p>
                  </li>
                ))}
              </ol>
            </div>

            <aside className="h-fit border-t-4 border-book-rust bg-secondary p-6 shadow-[4px_4px_0_var(--accent)] lg:mt-2">
              <h2 className="font-display text-xl font-semibold">Ready for another book?</h2>
              <p className="mt-3 text-sm leading-6 text-muted-foreground">
                Your catalog, current loans, and borrowing history are ready whenever you are.
              </p>
              <Link
                href={auth.status === "authenticated" ? "/dashboard" : "/login"}
                className="mt-4 inline-flex items-center gap-1 text-sm font-semibold text-primary underline-offset-4 hover:underline"
              >
                {auth.status === "authenticated" ? "Open my reading desk" : "See my loans"}
                <ArrowRight aria-hidden="true" className="size-3.5" />
              </Link>
            </aside>
          </div>
        </section>
      </main>
      <SiteFooter />
    </div>
  );
}
