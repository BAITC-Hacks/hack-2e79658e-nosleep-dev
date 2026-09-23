import type { Metadata } from "next";

import { LandingPage } from "@/components/landing-page";

export const metadata: Metadata = {
  title: "Mayor for five hours — city budget simulator",
  description: "Make five budget decisions and see how a city changes before you commit.",
};

export default function Home() {
  return <LandingPage />;
}
