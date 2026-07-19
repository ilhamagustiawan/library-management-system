import "@fontsource/dm-sans/400.css";
import "@fontsource/dm-sans/500.css";
import "@fontsource/dm-sans/600.css";
import "@fontsource/dm-sans/700.css";
import "@fontsource/fraunces/600.css";
import "@fontsource/fraunces/600-italic.css";
import "@fontsource/fraunces/700.css";
import "./globals.css";

import type { Metadata } from "next";

import { Toaster } from "@/components/ui/sonner";

import { Providers } from "./providers";

export const metadata: Metadata = {
  title: {
    default: "Perpus Digital — Your library, online",
    template: "%s — Perpus Digital",
  },
  description: "Browse available library books, borrow online, and keep every loan in view.",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <Providers>{children}</Providers>
        <Toaster position="bottom-right" />
      </body>
    </html>
  );
}
