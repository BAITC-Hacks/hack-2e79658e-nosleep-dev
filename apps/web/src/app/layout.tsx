import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Mayor for five hours — city budget simulator",
  description: "Five decisions, one budget, and a live view of district quality of life.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en" data-scroll-behavior="smooth">
      <body>{children}</body>
    </html>
  );
}
