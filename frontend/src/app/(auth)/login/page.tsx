"use client";

import { Info } from "lucide-react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AuthSession } from "@/features/auth/auth-session";
import { LoginForm } from "@/features/auth/login-form";

export default function LoginPage() {
  const router = useRouter();

  return (
    <Card className="bg-card/80 shadow-[6px_7px_0_0_var(--secondary)]">
      <CardHeader>
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Welcome back</p>
        <CardTitle as="h1">Open your member account</CardTitle>
        <CardDescription>
          Use any valid email and password. This MVP demonstrates the member journey with mock
          authentication.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <div className="mb-5 flex gap-2 rounded-md bg-secondary/70 p-3 text-xs leading-5 text-muted-foreground">
          <Info aria-hidden="true" className="mt-0.5 size-4 shrink-0 text-primary" />
          No credentials leave this browser. Passwords are never stored.
        </div>
        <LoginForm
          onAuthenticated={(session) => {
            AuthSession.write(window.localStorage, session);
            router.replace("/dashboard");
          }}
        />
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
