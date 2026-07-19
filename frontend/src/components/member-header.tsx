"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { BrandMark } from "@/components/brand-mark";
import { LogoutForm } from "@/features/auth/logout-form";
import { cn } from "@/lib/utils";

export function MemberHeader({ logoutEndpoint }: { logoutEndpoint: string }) {
  const pathname = usePathname();
  const dashboardCurrent = pathname === "/dashboard";
  const browseCurrent = pathname === "/browse" || pathname.startsWith("/books/");

  return (
    <header className="border-b border-border bg-card text-foreground">
      <div className="mx-auto flex min-h-16 max-w-7xl items-center gap-4 px-5 py-2 md:px-8">
        <Link
          aria-label="Perpus Digital"
          href="/dashboard?tab=books"
          className="inline-flex items-center gap-2 rounded-sm font-display text-2xl font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
        >
          <BrandMark nameClassName="hidden sm:inline" variant="on-light" />
        </Link>
        <nav aria-label="Member" className="ml-auto flex items-center gap-1">
          <Link
            aria-current={dashboardCurrent ? "page" : undefined}
            className={cn(
              "rounded-md px-2 py-2 text-sm font-semibold hover:bg-secondary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:px-3",
              dashboardCurrent && "bg-secondary text-primary",
            )}
            href="/dashboard?tab=books"
          >
            My Books
          </Link>
          <Link
            aria-current={browseCurrent ? "page" : undefined}
            className={cn(
              "rounded-md px-2 py-2 text-sm font-semibold hover:bg-secondary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring sm:px-3",
              browseCurrent && "bg-secondary text-primary",
            )}
            href="/browse"
          >
            Browse
          </Link>
        </nav>
        <LogoutForm logoutEndpoint={logoutEndpoint} tone="light" />
      </div>
    </header>
  );
}
