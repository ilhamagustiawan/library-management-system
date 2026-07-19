"use client";

import { LogOut } from "lucide-react";
import { useState, type FormEvent } from "react";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";

import { AuthApi } from "./auth-api";

export function LogoutForm({ logoutEndpoint }: { logoutEndpoint: string }) {
  const [state, setState] = useState<
    { status: "idle" } | { status: "pending" } | { status: "error"; message: string }
  >({ status: "idle" });

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const form = event.currentTarget;
    setState({ status: "pending" });
    const result = await AuthApi.logout(logoutEndpoint);
    if (result.status === "error") {
      setState(result);
      return;
    }
    form.submit();
  }

  return (
    <div>
      {state.status === "error" && <Alert className="mb-3">{state.message}</Alert>}
      <form action="/api/auth/logout" method="post" onSubmit={handleSubmit}>
        <Button variant="ghost" size="sm" type="submit" disabled={state.status === "pending"}>
          <LogOut aria-hidden="true" className="size-4" />
          {state.status === "pending" ? "Logging out…" : "Log out"}
        </Button>
      </form>
    </div>
  );
}
