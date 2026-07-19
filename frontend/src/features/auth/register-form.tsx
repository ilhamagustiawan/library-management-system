"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useMutation } from "@tanstack/react-query";
import { ArrowRight } from "lucide-react";
import { Controller, useForm } from "react-hook-form";

import { Alert } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

import { RegisterInput } from "./auth-schema";
import type { AuthSession } from "./auth-session";
import { FormError } from "./form-error";
import { MockAuth } from "./mock-auth";

export function RegisterForm({
  onAuthenticated,
}: {
  onAuthenticated: (session: AuthSession) => void;
}) {
  const form = useForm<RegisterInput>({
    resolver: zodResolver(RegisterInput.schema),
    defaultValues: {
      name: "",
      email: "",
      password: "",
      confirmPassword: "",
      acceptsTerms: false,
    },
  });
  const mutation = useMutation({
    mutationFn: MockAuth.register,
    onSuccess: (result) => {
      if (result.status === "success") {
        onAuthenticated(result.session);
        return;
      }

      form.setError("root", { message: result.error.message });
    },
  });

  return (
    <form className="space-y-4" noValidate onSubmit={form.handleSubmit((input) => mutation.mutate(input))}>
      {form.formState.errors.root?.message !== undefined && (
        <Alert>{form.formState.errors.root.message}</Alert>
      )}
      <div className="space-y-2">
        <Label htmlFor="name">Full name</Label>
        <Input
          id="name"
          autoComplete="name"
          aria-invalid={form.formState.errors.name !== undefined}
          aria-describedby={form.formState.errors.name ? "name-error" : undefined}
          placeholder="Maya Chen"
          {...form.register("name")}
        />
        <FormError id="name-error" message={form.formState.errors.name?.message} />
      </div>
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
      <div className="grid gap-4 sm:grid-cols-2">
        <div className="space-y-2">
          <Label htmlFor="password">Password</Label>
          <Input
            id="password"
            type="password"
            autoComplete="new-password"
            aria-invalid={form.formState.errors.password !== undefined}
            aria-describedby={form.formState.errors.password ? "password-error" : undefined}
            {...form.register("password")}
          />
          <FormError id="password-error" message={form.formState.errors.password?.message} />
        </div>
        <div className="space-y-2">
          <Label htmlFor="confirm-password">Confirm password</Label>
          <Input
            id="confirm-password"
            type="password"
            autoComplete="new-password"
            aria-invalid={form.formState.errors.confirmPassword !== undefined}
            aria-describedby={
              form.formState.errors.confirmPassword ? "confirm-password-error" : undefined
            }
            {...form.register("confirmPassword")}
          />
          <FormError
            id="confirm-password-error"
            message={form.formState.errors.confirmPassword?.message}
          />
        </div>
      </div>
      <Controller
        control={form.control}
        name="acceptsTerms"
        render={({ field }) => (
          <div className="space-y-2">
            <div className="flex items-start gap-3">
              <Checkbox
                id="terms"
                checked={field.value}
                onBlur={field.onBlur}
                onCheckedChange={(checked) => field.onChange(checked === true)}
                aria-invalid={form.formState.errors.acceptsTerms !== undefined}
                aria-describedby={
                  form.formState.errors.acceptsTerms ? "terms-error" : undefined
                }
              />
              <Label htmlFor="terms" className="font-normal leading-5 text-muted-foreground">
                I agree to the member terms and acknowledge this MVP uses mock authentication.
              </Label>
            </div>
            <FormError id="terms-error" message={form.formState.errors.acceptsTerms?.message} />
          </div>
        )}
      />
      <Button className="w-full" type="submit" disabled={mutation.isPending}>
        {mutation.isPending ? "Creating your account…" : "Create account"}
        {!mutation.isPending && <ArrowRight aria-hidden="true" className="size-4" />}
      </Button>
    </form>
  );
}
