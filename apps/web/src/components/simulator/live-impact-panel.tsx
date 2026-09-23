import { Activity, LoaderCircle, TrendingDown, TrendingUp } from "lucide-react";

import { CityMap } from "@/components/city-map";
import type { SimulationResponse } from "@/lib/types";

export function LiveImpactPanel({ simulation, simulating }: { simulation: SimulationResponse; simulating: boolean }) {
  const delta = simulation.dAvgAfter - simulation.dAvgBefore;
  const positive = delta >= 0;

  return (
    <section className="impact-panel" aria-labelledby="impact-title">
      <div className="impact-heading">
        <div><span className="impact-live"><i /> LIVE</span><h2 id="impact-title">Влияние на город</h2></div>
        {simulating ? <LoaderCircle className="spin" size={18} /> : <Activity size={18} />}
      </div>

      <div className="impact-score-row">
        <div className="impact-score"><span>Средний балл</span><strong>{simulation.dAvgAfter.toFixed(1)}</strong></div>
        <div className={`impact-delta ${positive ? "is-positive" : "is-negative"}`}>
          {positive ? <TrendingUp size={18} /> : <TrendingDown size={18} />}
          <span>{positive ? "+" : ""}{delta.toFixed(2)}<small> к базовой линии</small></span>
        </div>
      </div>

      <CityMap districts={simulation.districts} />

      <div className="impact-metrics">
        <div><span>Критические</span><strong>{simulation.criticalBefore} <small>→</small> {simulation.criticalAfter}</strong></div>
        <div><span>Самый уязвимый</span><strong>{simulation.districts.find((district) => district.id === simulation.minDistrictId)?.name ?? simulation.minDistrictId}</strong></div>
        <div><span>Средний балл</span><strong>{simulation.dAvgAfter.toFixed(1)}</strong></div>
      </div>
    </section>
  );
}
