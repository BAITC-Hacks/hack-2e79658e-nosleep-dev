import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "BesSheshim — AI city budget simulator",
  description: "Five decisions. One city. Explore how budget choices change quality of life across Astana.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru" data-scroll-behavior="smooth">
      <body>{children}</body>
    </html>
  );
}
