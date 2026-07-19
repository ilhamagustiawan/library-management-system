import Link from "next/link";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { AuthConfig } from "@/features/auth/auth-config";
import { RegisterForm } from "@/features/auth/register-form";

export const dynamic = "force-dynamic";

export default function RegisterPage() {
  const config = AuthConfig.load();

  return (
    <Card className="bg-card/80 shadow-[6px_7px_0_0_var(--secondary)]">
      <CardHeader>
        <p className="text-xs font-semibold uppercase tracking-[0.2em] text-primary">Join Libry</p>
        <CardTitle as="h1">Create your member account</CardTitle>
        <CardDescription>
          Start with the essentials. Library catalog and borrowing tools will connect when their
          APIs are ready.
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
