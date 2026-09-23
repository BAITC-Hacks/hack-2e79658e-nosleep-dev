import { ArrowLeft, Check, Sparkles } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Progress } from "@/components/ui/progress";
import type { ExplainResponse, OptimumResponse, SimulationResponse } from "@/lib/types";

type ScenarioResultProps = {
  explanation: ExplainResponse;
  optimum: OptimumResponse | null;
  simulation: SimulationResponse;
  onEdit: () => void;
};

export function ScenarioResult({ explanation, optimum, simulation, onEdit }: ScenarioResultProps) {
  const benchmark = optimum && simulation.score ? Math.min(100, (simulation.score / optimum.bestScore) * 100) : null;

  return (
    <section className="cockpit-result" id="result" aria-labelledby="result-title">
      <div className="cockpit-result-shell">
        <div className="result-overview">
          <button type="button" onClick={onEdit}><ArrowLeft size={16} /> Изменить план</button>
          <div className="result-label"><Sparkles size={15} /> AI-РАЗБОР СЦЕНАРИЯ</div>
          <h2 id="result-title">Город стал<br /><span>сильнее.</span></h2>
          <p>{explanation.explanation.summary}</p>
          <Badge>{explanation.source === "cached" ? "Кэшированный разбор" : "Live AI"}</Badge>
        </div>

        <div className="result-score-card">
          <span>Итоговый индекс</span>
          <strong>{simulation.score?.toFixed(2)}</strong>
          <div><span>Средний балл</span><strong>{simulation.dAvgBefore.toFixed(1)} → {simulation.dAvgAfter.toFixed(1)}</strong></div>
          <div><span>Критические показатели</span><strong>{simulation.criticalBefore} → {simulation.criticalAfter}</strong></div>
          {benchmark !== null && <div className="result-optimum"><span>От лучшего сценария <strong>{benchmark.toFixed(1)}%</strong></span><Progress value={benchmark} /></div>}
        </div>

        <div className="cockpit-result-columns">
          <ResultColumn title="Что сработало" items={explanation.explanation.strengths} positive />
          <ResultColumn title="Риски" items={explanation.explanation.risks} />
          <ResultColumn title="Следующий шаг" items={explanation.explanation.recommendations} />
        </div>
      </div>
    </section>
  );
}

function ResultColumn({ title, items, positive = false }: { title: string; items: string[]; positive?: boolean }) {
  return <div><h3>{title}</h3>{items.map((item, index) => <p key={item}>{positive ? <Check size={16} /> : <span>0{index + 1}</span>}{item}</p>)}</div>;
}
