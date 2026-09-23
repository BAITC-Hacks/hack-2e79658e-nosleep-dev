import Link from "next/link";
import { CircleHelp, House, Layers3, RotateCcw } from "lucide-react";

import { Progress } from "@/components/ui/progress";

type ScenarioHeaderProps = {
  budget: number;
  budgetRemaining: number;
  decisionCount: number;
  requiredDecisions: number;
  health: "loading" | "ok" | "error";
  isMock: boolean;
  onReset: () => void;
  onHelp: () => void;
  savedCount: number;
  onCompare: () => void;
};

export function ScenarioHeader({ budget, budgetRemaining, decisionCount, requiredDecisions, health, isMock, onReset, onHelp, savedCount, onCompare }: ScenarioHeaderProps) {
  const budgetUsed = budget - budgetRemaining;

  return (
    <header className="cockpit-header">
      <Link className="cockpit-brand" href="/" aria-label="BesSheshim — на главную">
        <span className="cockpit-brand-mark">BS/5</span>
        <span>BES SHESHIM</span>
      </Link>

      <div className="cockpit-session">
        <span className="cockpit-kicker">СЦЕНАРИЙ 01</span>
        <div className="cockpit-session-progress">
          <strong>{decisionCount}<small> / {requiredDecisions}</small></strong>
          <span>решений принято</span>
        </div>
      </div>

      <div className="cockpit-budget">
        <div><span>Осталось</span><strong>{budgetRemaining}<small> / {budget}</small></strong></div>
        <Progress value={(budgetUsed / budget) * 100} />
      </div>

      <div className="cockpit-header-actions">
        <span className={`cockpit-api ${health}`}><i /> {isMock ? "DEMO" : health}</span>
        {savedCount > 0 && <button className="cockpit-icon-button" type="button" onClick={onCompare} aria-label={`Сравнить сценарии (${savedCount})`}><Layers3 size={18} /></button>}
        <button className="cockpit-icon-button" type="button" onClick={onHelp} aria-label="Как пользоваться симулятором"><CircleHelp size={18} /></button>
        <button className="cockpit-icon-button" type="button" onClick={onReset} disabled={decisionCount === 0} aria-label="Сбросить сценарий"><RotateCcw size={18} /></button>
        <Link className="cockpit-home" href="/"><House size={16} /> <span>Главная</span></Link>
      </div>
    </header>
  );
}
