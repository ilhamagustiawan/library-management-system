import { BookMarked, BookOpen, Sparkles } from "lucide-react";
import Link from "next/link";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";

import type { OAuthUser } from "./oauth-client";
import { LogoutForm } from "./logout-form";

export function DashboardClient({
  logoutEndpoint,
  session,
}: {
  logoutEndpoint: string;
  session: OAuthUser;
}) {
  return (
    <div className="min-h-screen">
      <header className="border-b border-border/70">
        <div className="mx-auto flex h-18 max-w-6xl items-center justify-between px-5 md:px-8">
          <Link href="/" className="inline-flex items-center gap-2 font-display text-2xl font-semibold">
            <span className="grid size-8 place-items-center rounded-full bg-primary text-primary-foreground">
              <BookOpen aria-hidden="true" className="size-4" />
            </span>
            Libry
          </Link>
          <LogoutForm logoutEndpoint={logoutEndpoint} />
        </div>
      </header>
      <main className="mx-auto max-w-6xl px-5 py-12 md:px-8 md:py-18">
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Member home</p>
        <h1 className="mt-3 font-display text-5xl font-semibold tracking-tight sm:text-6xl">
          Welcome, {session.name}.
        </h1>
        <p className="mt-4 max-w-xl text-lg leading-8 text-muted-foreground">
          Your quiet corner of the library is ready. More member tools will appear here as services
          connect.
        </p>

        <div className="mt-10 grid gap-5 lg:grid-cols-[0.8fr_1.2fr]">
          <Card className="overflow-hidden bg-primary text-primary-foreground">
            <CardHeader>
              <BookMarked aria-hidden="true" className="size-7 opacity-80" strokeWidth={1.6} />
              <CardTitle className="pt-8 text-primary-foreground">Member card</CardTitle>
              <CardDescription className="text-primary-foreground/70">
                Authenticated membership · active session
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="border-t border-primary-foreground/20 pt-5">
                <p className="text-sm font-semibold">{session.name}</p>
                <p className="mt-1 text-sm text-primary-foreground/70">{session.email}</p>
              </div>
            </CardContent>
          </Card>

          <Card className="border-dashed bg-card/55">
            <CardHeader>
              <Sparkles aria-hidden="true" className="size-7 text-primary" strokeWidth={1.6} />
              <CardTitle className="pt-8">Your reading activity will live here</CardTitle>
              <CardDescription className="max-w-lg">
                Catalog search, holds, loans, and due dates are intentionally omitted from this MVP
                until backend contracts exist.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>
      </main>
    </div>
  );
}
