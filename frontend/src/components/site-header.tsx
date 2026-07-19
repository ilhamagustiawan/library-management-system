import Link from "next/link";

import { BrandMark } from "@/components/brand-mark";
import { buttonVariants } from "@/components/ui/button";
import { LogoutForm } from "@/features/auth/logout-form";
import { cn } from "@/lib/utils";

export type SiteHeaderAuth =
  | { status: "guest" }
  | { status: "authenticated"; logoutEndpoint: string };

export function SiteHeader({ auth }: { auth: SiteHeaderAuth }) {
  return (
    <header className="border-b border-foreground bg-foreground text-background">
      <div className="mx-auto flex h-16 max-w-6xl items-center justify-between px-5 md:px-8">
        <Link
          href={auth.status === "authenticated" ? "/dashboard?tab=books" : "/"}
          className="inline-flex items-center gap-2 rounded-sm font-display text-2xl font-semibold focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-background"
        >
          <BrandMark variant="on-dark" />
        </Link>
        <nav aria-label="Primary" className="flex items-center gap-1 sm:gap-2">
          {auth.status === "authenticated" ? (
            <>
              <Link
                href="/dashboard"
                className={cn(
                  buttonVariants({ variant: "ghost", size: "sm" }),
                  "text-background hover:bg-background/10",
                )}
              >
                Member home
              </Link>
              <Link
                href="/browse"
                className={cn(
                  buttonVariants({ variant: "ghost", size: "sm" }),
                  "hidden text-background hover:bg-background/10 sm:inline-flex",
                )}
              >
                Browse
              </Link>
              <LogoutForm logoutEndpoint={auth.logoutEndpoint} />
            </>
          ) : (
            <>
              <Link
                href="/login"
                className={cn(
                  buttonVariants({ variant: "ghost", size: "sm" }),
                  "text-background hover:bg-background/10",
                )}
              >
                Log in
              </Link>
              <Link
                href="/register"
                className={cn(
                  buttonVariants({ variant: "outline", size: "sm" }),
                  "hidden border-background bg-background text-foreground hover:bg-secondary sm:inline-flex",
                )}
              >
                Join Perpus Digital
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
