"use client";

import { useRouter } from "next/navigation";
import { useEffect, useRef } from "react";
import { toast } from "sonner";

import type { DashboardNotice } from "./dashboard-notice";

export function DashboardNoticeToast({ notice }: { notice: DashboardNotice }) {
  const router = useRouter();
  const shown = useRef(false);

  useEffect(() => {
    if (notice.kind !== "book-borrowed" || shown.current) return;
    shown.current = true;
    toast.success("Book added to your shelf.", {
      description: "Your active loans are now up to date.",
    });
    const params = new URLSearchParams(window.location.search);
    params.delete("borrowed");
    const query = params.toString();
    router.replace(`/dashboard${query === "" ? "" : `?${query}`}`, { scroll: false });
  }, [notice, router]);

  return null;
}
