import type { Metadata } from "next";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { MemberHeader } from "@/components/member-header";
import { AuthConfig } from "@/features/auth/auth-config";
import { AuthCookies } from "@/features/auth/auth-cookies";
import { SessionRefresh } from "@/features/auth/session-refresh";
import { WebSession } from "@/features/auth/web-session";
import { BookBrowser } from "@/features/library/book-browser";
import { BookCatalog } from "@/features/library/book-catalog";

export const dynamic = "force-dynamic";

export const metadata: Metadata = { title: "Browse books" };

export default async function BrowsePage() {
  const config = AuthConfig.load();
  const sealedSession = (await cookies()).get(AuthCookies.sessionName)?.value;
  if (sealedSession === undefined) redirect("/login");
  const openedSession = WebSession.open(sealedSession, config.sessionSecret);
  if (openedSession.status === "invalid") redirect("/login");
  if (WebSession.needsRefresh(openedSession.session)) return <SessionRefresh />;

  const catalog = await BookCatalog.listAvailable({
    issuer: config.oauth.serviceURL,
    accessToken: openedSession.session.accessToken,
  });
  if (catalog.status === "error" && catalog.error.kind === "unauthorized") redirect("/login");

  return (
    <div className="flex min-h-screen flex-col">
      <MemberHeader logoutEndpoint={config.logoutEndpoint} />
      <main className="mx-auto w-full max-w-7xl flex-1 px-5 py-8 md:px-8 md:py-12">
        <p className="text-xs font-bold uppercase tracking-[0.2em] text-primary">The stacks</p>
        <h1 className="mt-2 font-display text-5xl font-semibold leading-none tracking-[-0.03em] sm:text-6xl">Browse available books</h1>
        <p className="mt-3 max-w-2xl text-sm leading-6 text-muted-foreground sm:text-base">
          Every title below has a copy ready to borrow. Select a book for full details and availability.
        </p>
        {catalog.status === "success" ? (
          <BookBrowser books={catalog.books} />
        ) : (
          <div role="alert" className="mt-8 border-l-4 border-destructive bg-card p-5">
            <p className="font-semibold">Catalog unavailable</p>
            <p className="mt-1 text-sm text-muted-foreground">Books could not load. Your account remains unchanged. Try again.</p>
          </div>
        )}
      </main>
      <footer className="border-t border-border bg-secondary px-5 py-5 text-center text-xs text-muted-foreground">
        Perpus Digital member portal
      </footer>
    </div>
  );
}
