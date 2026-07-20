"use client";

import { LogOut } from "lucide-react";
import { useState, type FormEvent } from "react";

import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

import { AuthApi } from "./auth-api";

export function LogoutForm({
  logoutEndpoint,
  tone = "dark",
}: {
  logoutEndpoint: string;
  tone?: "dark" | "light";
}) {
  const [state, setState] = useState<{ status: "idle" } | { status: "pending" }>({
    status: "idle",
  });

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    setState({ status: "pending" });
    // Local teardown must not depend on Auth Service availability.
    void AuthApi.logout(logoutEndpoint);
    form.submit();
  }

  return (
    <div>
      <form action="/api/auth/logout" method="post" onSubmit={handleSubmit}>
        <Button
          className={cn(
            tone === "dark"
              ? "text-background hover:bg-background/10"
              : "text-foreground hover:bg-secondary",
          )}
          variant="ghost"
          size="sm"
          type="submit"
          disabled={state.status === "pending"}
        >
          <LogOut aria-hidden="true" className="size-4" />
          {state.status === "pending" ? "Logging out…" : "Log out"}
        </Button>
      </form>
    </div>
  );
}
