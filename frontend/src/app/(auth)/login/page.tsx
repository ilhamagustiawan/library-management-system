import { Info } from "lucide-react";
import Link from "next/link";

import { Alert } from "@/components/ui/alert";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AuthConfig } from "@/features/auth/auth-config";
import { LoginForm } from "@/features/auth/login-form";
import { OAuthNavigation } from "@/features/auth/oauth-navigation";
import { OAuthStart } from "@/features/auth/oauth-start";

export const dynamic = "force-dynamic";

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
    <Card className="bg-card/80 shadow-[6px_7px_0_0_var(--secondary)]">
      <CardHeader>
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Welcome back</p>
        <CardTitle as="h1">Open your member account</CardTitle>
        <CardDescription>
          Sign in through the library authorization service. Access and refresh tokens stay in an
          encrypted, HttpOnly server session.
        </CardDescription>
      </CardHeader>
      <CardContent>
        {error !== undefined && (
          <Alert>{errors[error] ?? "Login failed safely. Start again; no session was created."}</Alert>
        )}
        <div className="mb-5 flex gap-2 rounded-md bg-secondary/70 p-3 text-xs leading-5 text-muted-foreground">
          <Info aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-primary" />
          Credentials go directly to the auth service. Passwords and tokens are never stored in
          browser-accessible storage.
        </div>
        {validReturnTo ? (
          <LoginForm loginEndpoint={config.loginEndpoint} returnTo={returnTo} />
        ) : (
          <OAuthStart automatic={error === undefined} />
        )}
        <p className="mt-6 text-center text-sm text-muted-foreground">
          New to Libry?{" "}
          <Link className="font-semibold text-foreground underline-offset-4 hover:underline" href="/register">
            Create an account
          </Link>
        </p>
      </CardContent>
    </Card>
  );
}
