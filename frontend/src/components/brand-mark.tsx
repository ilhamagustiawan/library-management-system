import { BookOpen } from "lucide-react";

import { cn } from "@/lib/utils";

type BrandMarkProps = {
  nameClassName?: string;
  variant: "on-dark" | "on-light";
};

export function BrandMark({ nameClassName, variant }: BrandMarkProps) {
  const onDark = variant === "on-dark";

  return (
    <>
      <span
        aria-hidden="true"
        className={cn(
          "grid size-8 shrink-0 place-items-center rounded-sm shadow-[3px_3px_0_var(--book-rust)]",
          onDark ? "bg-background text-foreground" : "bg-foreground text-background",
        )}
      >
        <BookOpen className="size-4" strokeWidth={1.8} />
      </span>
      <span className={nameClassName}>
        Perpus{" "}
        <span className={cn("italic", onDark ? "text-book-gold" : "text-book-rust")}>Digital</span>
      </span>
    </>
  );
}
