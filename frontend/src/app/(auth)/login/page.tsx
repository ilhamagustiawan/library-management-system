import { Info } from "lucide-react";
import type { Metadata } from "next";
import Link from "next/link";

import { Alert } from "@/components/ui/alert";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AuthConfig } from "@/features/auth/auth-config";
import { LoginForm } from "@/features/auth/login-form";
import { OAuthNavigation } from "@/features/auth/oauth-navigation";
import { OAuthStart } from "@/features/auth/oauth-start";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Log in",
  description: "Log in to browse books, borrow online, and review your Perpus Digital loans.",
};

type LoginPageProps = {
  searchParams: Promise<{ error?: string | string[]; return_to?: string | string[] }>;
};

const errors: Record<string, string> = {
  authorization_denied: "Authorization was cancelled. No session was created.",
  invalid_callback: "Login response could not be verified. Start again to protect your account.",
  session_expired: "Your session expired and could not be refreshed. Log in again.",
  token_exchange_failed: "Login code could not be exchanged. Start again; no session was created.",
  user_info_failed: "Your account could not be loaded. Start login again.",
};

function first(value: string | string[] | undefined) {
  return Array.isArray(value) ? value[0] : value;
}

export default async function LoginPage({ searchParams }: LoginPageProps) {
  const config = AuthConfig.load();
  const params = await searchParams;
  const error = first(params.error);
  const returnTo = first(params.return_to);
  const validReturnTo =
    returnTo !== undefined && OAuthNavigation.isAuthorizeReturnTo(returnTo, config.oauth.issuer);

  return (
    <Card className="border-t-4 border-t-book-rust shadow-[6px_6px_0_var(--accent)]">
      <CardHeader>
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Welcome back</p>
        <CardTitle as="h1" className="text-3xl sm:text-4xl">Return to your reading</CardTitle>
        <CardDescription>
          Log in to browse available books, borrow your next read, and check every loan in one
          place.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {error !== undefined && (
          <Alert>{errors[error] ?? "Login failed safely. Start again; no session was created."}</Alert>
        )}
        <div className="mb-5 flex gap-2 rounded-sm border border-border bg-secondary p-3 text-xs leading-5 text-muted-foreground">
          <Info aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-primary" />
          Your sign-in is protected. Perpus Digital keeps your password and session details out of
          browser storage.
        </div>
        {validReturnTo ? (
          <LoginForm loginEndpoint={config.loginEndpoint} returnTo={returnTo} />
        ) : (
          <OAuthStart automatic={error === undefined} />
        )}
        <p className="mt-6 text-center text-sm text-muted-foreground">
          New to Perpus Digital?{" "}
          <Link className="font-semibold text-foreground underline-offset-4 hover:underline" href="/register">
            Create an account
          </Link>
        </p>
      </CardContent>
    </Card>
  );
}
