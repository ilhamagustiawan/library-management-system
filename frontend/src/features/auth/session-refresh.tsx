"use client";

import { useEffect, useRef } from "react";

type RunExclusive = (operation: () => Promise<void>) => Promise<void>;

function browserNavigate(url: string) {
  window.location.replace(url);
}

async function browserExclusive(operation: () => Promise<void>) {
  if (navigator.locks === undefined) {
    await operation();
    return;
  }
  await navigator.locks.request("lms-session-refresh", operation);
}

export function SessionRefresh({
  fetcher = fetch,
  navigate = browserNavigate,
  runExclusive = browserExclusive,
}: {
  fetcher?: typeof fetch;
  navigate?: (url: string) => void;
  runExclusive?: RunExclusive;
}) {
  const started = useRef(false);

  useEffect(() => {
    if (started.current) return;
    started.current = true;
    void runExclusive(async () => {
      try {
        const response = await fetcher("/api/auth/refresh", {
          method: "POST",
          headers: { Accept: "application/json" },
        });
        navigate(response.ok ? "/dashboard" : "/login?error=session_expired");
      } catch {
        navigate("/login?error=session_expired");
      }
    });
  }, [fetcher, navigate, runExclusive]);

  return (
    <main className="grid min-h-screen place-items-center bg-secondary px-5" aria-busy="true">
      <div className="border border-border bg-card px-6 py-5 text-center">
        <p className="font-display text-lg font-semibold text-foreground">Perpus Digital</p>
        <p className="mt-2 text-sm text-muted-foreground">Refreshing secure session…</p>
      </div>
    </main>
  );
}
