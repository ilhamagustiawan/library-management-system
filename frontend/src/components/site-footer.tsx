import Link from "next/link";

export function SiteFooter() {
  return (
    <footer className="border-t border-border/70">
      <div className="mx-auto flex max-w-7xl flex-col gap-3 px-5 py-8 text-sm text-muted-foreground sm:flex-row sm:items-center sm:justify-between md:px-8">
        <p>Libry member portal · MVP preview</p>
        <nav aria-label="Footer" className="flex gap-5">
          <Link className="underline-offset-4 hover:text-foreground hover:underline" href="/login">
            Log in
          </Link>
          <Link
            className="underline-offset-4 hover:text-foreground hover:underline"
            href="/register"
          >
            Register
          </Link>
        </nav>
      </div>
    </footer>
  );
}
