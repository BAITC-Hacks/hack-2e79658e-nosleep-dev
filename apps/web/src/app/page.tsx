import type { Metadata } from "next";

import { LandingPage } from "@/components/landing-page";

export const metadata: Metadata = {
  title: "BesSheshim — AI city budget simulator",
  description: "Решения, из которых строится город. Управляйте бюджетом Астаны и сразу увидьте последствия каждого выбора.",
};

export default function Home() {
  return <LandingPage />;
}
