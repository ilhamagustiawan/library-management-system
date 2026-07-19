"use client";

import { cn } from "@/lib/utils";

const palettes = [
  "bg-primary text-primary-foreground",
  "bg-book-rust text-primary-foreground",
  "bg-book-gold text-foreground",
  "bg-foreground text-background",
] as const;

type BookCoverProps = {
  coverUrl: string | null;
  size?: "thumbnail" | "loan" | "card" | "detail";
  title: string;
};

export function BookCover({ coverUrl, title, size = "card" }: BookCoverProps) {
  const palette = palettes[title.length % palettes.length];

  return (
    <div
      aria-hidden="true"
      data-size={size}
      className={cn(
        palette,
        "relative flex shrink-0 flex-col justify-between overflow-hidden border border-foreground/20 shadow-[0_10px_24px_rgba(33,26,20,0.16)]",
        size === "thumbnail" && "h-20 w-14 p-2",
        size === "loan" && "aspect-[2/3] w-20 p-3 sm:w-24",
        size === "card" && "aspect-[2/3] w-full p-4",
        size === "detail" && "aspect-[2/3] w-full max-w-64 p-4 sm:w-64",
      )}
    >
      <span className="h-px w-10 bg-current opacity-60" />
      <span
        className={cn(
          "font-display font-semibold leading-tight",
          size === "thumbnail"
            ? "text-xl uppercase"
            : size === "loan"
              ? "text-base sm:text-lg"
              : "text-xl sm:text-2xl",
        )}
      >
        {size === "thumbnail" ? title.charAt(0) : title}
      </span>
      {size !== "thumbnail" && (
        <span className="text-[0.65rem] font-bold uppercase tracking-[0.2em] opacity-75">
          Perpus Digital edition
        </span>
      )}
      {coverUrl !== null && (
        // Catalog cover hosts are librarian-controlled, so a native image preserves valid
        // HTTP(S) sources without weakening Next.js host allow-listing.
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={coverUrl}
          alt=""
          width="400"
          height="600"
          decoding="async"
          loading={size === "detail" ? "eager" : "lazy"}
          fetchPriority={size === "detail" ? "high" : "auto"}
          className="absolute inset-0 size-full bg-card object-contain opacity-0 transition-opacity duration-[var(--motion-duration-ui)] ease-[var(--motion-ease-out)] data-[loaded=true]:opacity-100"
          onLoad={(event) => {
            event.currentTarget.dataset.loaded = "true";
          }}
          onError={(event) => {
            event.currentTarget.hidden = true;
          }}
        />
      )}
    </div>
  );
}
