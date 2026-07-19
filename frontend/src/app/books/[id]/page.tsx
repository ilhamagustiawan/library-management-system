import { ArrowLeft, BookCopy, CalendarDays, Hash, UserRound } from "lucide-react";
import type { Metadata } from "next";
import { cookies } from "next/headers";
import Link from "next/link";
import { notFound, redirect } from "next/navigation";
import { z } from "zod";

import { MemberHeader } from "@/components/member-header";
import { Button } from "@/components/ui/button";
import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { SessionRefresh } from "@/features/auth/session-refresh";
import { WebSession } from "@/features/auth/web-session";
import { BookCatalog } from "@/features/library/book-catalog";
import { BookCover } from "@/features/library/book-cover";
import { BorrowButton } from "@/features/library/borrow-button";
import { LoanEligibility } from "@/features/library/loan-eligibility";
import { MemberLibrary } from "@/features/library/member-library";
import { ReturnBookDialog } from "@/features/library/return-book-dialog";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Book details" };

export default async function BookDetailPage({ params }: { params: Promise<{ id: string }> }) {
  const { id } = await params;
  if (!z.uuid().safeParse(id).success) notFound();
  const config = AuthConfig.load();
  const sealedSession = (await cookies()).get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) redirect("/login");
  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid") redirect("/login");
  if (WebSession.needsRefresh(openedSession.session)) return <SessionRefresh />;

  const [catalog, library] = await Promise.all([
    BookCatalog.get({ issuer: config.oauth.serviceURL, accessToken: openedSession.session.accessToken, id }),
    MemberLibrary.load({ issuer: config.oauth.serviceURL, accessToken: openedSession.session.accessToken }),
  ]);
  if (catalog.status === "error" && catalog.error.kind === "unauthorized") redirect("/login");
  if (catalog.status === "error" && catalog.error.kind === "not-found") notFound();
  const activeLoan = catalog.status === "success"
    ? MemberLibrary.activeLoanForBook(library, catalog.book.id)
    : undefined;

  return (
    <div className="flex min-h-screen flex-col">
      <MemberHeader logoutEndpoint={config.logoutEndpoint} />
      <main className="mx-auto w-full max-w-7xl flex-1 px-5 py-8 md:px-8 md:py-12">
        <Link href="/browse" className="inline-flex items-center gap-2 text-sm font-semibold text-primary underline-offset-4 hover:underline">
          <ArrowLeft aria-hidden="true" className="size-4" /> Back to Browse
        </Link>
        {catalog.status === "error" ? (
          <div role="alert" className="mt-8 border-l-4 border-destructive bg-card p-5">
            <h1 className="font-display text-2xl font-semibold">Book unavailable</h1>
            <p className="mt-2 text-sm text-muted-foreground">Details could not load. Your account remains unchanged. Try again.</p>
          </div>
        ) : (
          <article className="mt-7 grid gap-8 md:grid-cols-[16rem_minmax(0,1fr)] lg:gap-12">
            <aside>
              <BookCover
                coverUrl={catalog.book.coverUrl}
                title={catalog.book.title}
                size="detail"
              />
              <div className="mt-5 border-y border-border py-5">
                {activeLoan === undefined ? (
                  <BorrowButton bookId={catalog.book.id} eligibility={LoanEligibility.forBook(catalog.book, library)} />
                ) : (
                  <div>
                    <p className="mb-3 text-sm font-semibold text-primary">This book is on your shelf.</p>
                    <ReturnBookDialog bookTitle={catalog.book.title} loanId={activeLoan.loanId}>
                      <Button className="h-12 w-full text-base" variant="outline">Return this book</Button>
                    </ReturnBookDialog>
                  </div>
                )}
              </div>
            </aside>
            <div className="min-w-0">
              <p className="text-xs font-bold uppercase tracking-[0.16em] text-primary">Available in Perpus Digital</p>
              <h1 className="mt-2 max-w-4xl font-display text-4xl font-semibold leading-[1.05] tracking-tight sm:text-5xl lg:text-6xl">
                {catalog.book.title}
              </h1>
              <p className="mt-4 font-display text-2xl text-muted-foreground">by {catalog.book.author}</p>
              <dl className="mt-8 grid gap-px border border-border bg-border sm:grid-cols-3">
                <div className="bg-card p-4">
                  <dt className="flex items-center gap-2 text-xs font-bold uppercase tracking-[0.12em] text-muted-foreground"><BookCopy aria-hidden="true" className="size-4" /> Copies</dt>
                  <dd className="mt-2 font-display text-xl font-semibold">{catalog.book.availableCopies} of {catalog.book.totalCopies} ready</dd>
                </div>
                <div className="bg-card p-4">
                  <dt className="flex items-center gap-2 text-xs font-bold uppercase tracking-[0.12em] text-muted-foreground"><CalendarDays aria-hidden="true" className="size-4" /> Published</dt>
                  <dd className="mt-2 font-display text-xl font-semibold">{catalog.book.publicationYear ?? "Unknown"}</dd>
                </div>
                <div className="bg-card p-4">
                  <dt className="flex items-center gap-2 text-xs font-bold uppercase tracking-[0.12em] text-muted-foreground"><Hash aria-hidden="true" className="size-4" /> ISBN</dt>
                  <dd className="mt-2 break-all font-display text-xl font-semibold">{catalog.book.isbn}</dd>
                </div>
              </dl>
              <section className="mt-9 border-t border-border pt-7">
                <h2 className="font-display text-3xl font-semibold">About this book</h2>
                <p className="mt-4 max-w-3xl whitespace-pre-line text-base leading-8 text-card-foreground">
                  {catalog.book.description ?? "No description available for this edition."}
                </p>
              </section>
              <section className="mt-9 border-t border-border pt-7">
                <h2 className="flex items-center gap-3 font-display text-2xl font-semibold"><UserRound aria-hidden="true" className="size-5 text-primary" /> Author</h2>
                <p className="mt-3 text-lg">{catalog.book.author}</p>
              </section>
            </div>
          </article>
        )}
      </main>
      <footer className="border-t border-border bg-secondary px-5 py-5 text-center text-xs text-muted-foreground">Perpus Digital member portal</footer>
    </div>
  );
}
