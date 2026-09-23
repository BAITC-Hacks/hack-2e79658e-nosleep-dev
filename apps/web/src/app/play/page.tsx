import type { Metadata } from "next";

import { Simulator } from "@/components/simulator";

export const metadata: Metadata = {
  title: "Симулятор — BesSheshim",
  description: "Соберите пять городских решений и увидьте их влияние на районы Астаны.",
};

export default function PlayPage() {
  return <Simulator />;
}
