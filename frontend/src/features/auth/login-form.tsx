"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { ArrowRight } from "lucide-react";
import { useForm } from "react-hook-form";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { LoginInput } from "./auth-schema";
import type { AuthSession } from "./auth-session";
import { FormError } from "./form-error";
import { MockAuth } from "./mock-auth";

export function LoginForm({
  onAuthenticated,
}: {
  onAuthenticated: (session: AuthSession) => void;
}) {
  const form = useForm<LoginInput>({
    resolver: zodResolver(LoginInput.schema),
    defaultValues: { email: "", password: "" },
  });
  const mutation = useMutation({
    mutationFn: MockAuth.login,
    onSuccess: (result) => {
      if (result.status === "success") {
        onAuthenticated(result.session);
        return;
      }

      form.setError("root", { message: result.error.message });
    },
  });

  return (
    <form className="space-y-5" noValidate onSubmit={form.handleSubmit((input) => mutation.mutate(input))}>
      {form.formState.errors.root?.message !== undefined && (
        <Alert>{form.formState.errors.root.message}</Alert>
      )}
      <div className="space-y-2">
        <Label htmlFor="email">Email address</Label>
        <Input
          id="email"
          type="email"
          autoComplete="email"
          aria-invalid={form.formState.errors.email !== undefined}
          aria-describedby={form.formState.errors.email ? "email-error" : undefined}
          placeholder="you@example.com"
          {...form.register("email")}
        />
        <FormError id="email-error" message={form.formState.errors.email?.message} />
      </div>
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <Label htmlFor="password">Password</Label>
          <span className="text-xs text-muted-foreground">8+ characters</span>
        </div>
        <Input
          id="password"
          type="password"
          autoComplete="current-password"
          aria-invalid={form.formState.errors.password !== undefined}
          aria-describedby={form.formState.errors.password ? "password-error" : undefined}
          {...form.register("password")}
        />
        <FormError id="password-error" message={form.formState.errors.password?.message} />
      </div>
      <Button className="w-full" type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? "Opening your account…" : "Log in"}
        {!mutation.isPending && <ArrowRight aria-hidden="true" className="size-4" />}
      </Button>
    </form>
  );
}
