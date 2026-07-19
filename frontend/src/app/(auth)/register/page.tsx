import type { Metadata } from "next";
import Link from "next/link";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AuthConfig } from "@/features/auth/auth-config";
import { RegisterForm } from "@/features/auth/register-form";

export const dynamic = "force-dynamic";

export const metadata: Metadata = {
  title: "Create account",
  description:
    "Create a Perpus Digital account to browse books, borrow online, and track every loan.",
};

export default function RegisterPage() {
  const config = AuthConfig.load();

  return (
    <Card className="border-t-4 border-t-book-rust shadow-[6px_6px_0_var(--accent)]">
      <CardHeader>
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">
          Join Perpus Digital
        </p>
        <CardTitle as="h1" className="text-3xl sm:text-4xl">
          Create your Perpus Digital account
        </CardTitle>
        <CardDescription>
          Create your member account to browse available books and keep your loans and reading
          history together.
        </CardDescription>
      </CardHeader>
      <CardContent>
        <RegisterForm registerEndpoint={config.registerEndpoint} />
        <p className="mt-6 text-center text-sm text-muted-foreground">
          Already registered?{" "}
          <Link className="font-semibold text-foreground underline-offset-4 hover:underline" href="/login">
            Log in
          </Link>
        </p>
      </CardContent>
    </Card>
  );
}
