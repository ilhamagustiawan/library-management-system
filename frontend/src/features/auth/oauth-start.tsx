"use client";

import { useEffect } from "react";

import { buttonVariants } from "@/components/ui/button";

const endpoint = "/api/auth/login";

export function OAuthStart({
  automatic = true,
  navigate,
}: {
  automatic?: boolean;
  navigate?: (url: string) => void;
}) {
  useEffect(() => {
    if (!automatic) return;
    if (navigate !== undefined) navigate(endpoint);
    else window.location.assign(endpoint);
  }, [automatic, navigate]);

  return (
    <a className={buttonVariants({ className: "w-full" })} href={endpoint}>
      {automatic ? "Continue to login" : "Start login again"}
    </a>
  );
}
