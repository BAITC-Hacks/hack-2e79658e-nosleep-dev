import type { Metadata } from "next";

import { Simulator } from "@/components/simulator";

export const metadata: Metadata = {
  title: "Play — Mayor for five hours",
};

export default function PlayPage() {
  return <Simulator />;
}
