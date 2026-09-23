import { ArrowRight, CircleAlert, LoaderCircle, MapPin, RotateCcw, X } from "lucide-react";

import { Button } from "@/components/ui/button";
import type { CatalogResponse, Decision, SimulationResponse } from "@/lib/types";

type PlanTrayProps = {
  catalog: CatalogResponse;
  decisions: Decision[];
  error: string | null;
  explaining: boolean;
  simulating: boolean;
  simulation: SimulationResponse;
  onRemove: (initiativeId: string) => void;
  onReset: () => void;
  onRetry: () => void;
  onSubmit: () => void;
};

export function PlanTray({ catalog, decisions, error, explaining, simulating, simulation, onRemove, onReset, onRetry, onSubmit }: PlanTrayProps) {
  const remaining = catalog.requiredDecisions - decisions.length;
  const cta = simulation.submittable ? "Посмотреть результат" : remaining > 0 ? `Выберите ещё ${remaining}` : "Исправьте ограничения";

  return (
    <aside className="cockpit-panel plan-tray" aria-labelledby="plan-title">
      <div className="cockpit-panel-heading">
        <div><span className="cockpit-kicker">ШАГ 02</span><h2 id="plan-title">Ваш план</h2></div>
        <button className="plan-reset" type="button" onClick={onReset} disabled={decisions.length === 0}><RotateCcw size={14} /> Сбросить</button>
      </div>

      <div className="plan-slots">
        {Array.from({ length: catalog.requiredDecisions }, (_, index) => {
          const decision = decisions[index];
          const initiative = decision ? catalog.initiatives.find((item) => item.id === decision.initiativeId) : undefined;
          const district = decision?.districtId ? catalog.districts.find((item) => item.id === decision.districtId) : undefined;
          return (
            <div className={`plan-slot ${initiative ? "is-filled" : ""}`} key={index}>
              <span className="plan-slot-index">0{index + 1}</span>
              {initiative ? <>
                <div><strong>{initiative.name}</strong><span>{initiative.cost} ед.{district && <> · <MapPin size={11} /> {district.name}</>}</span></div>
                <button type="button" onClick={() => onRemove(initiative.id)} aria-label={`Убрать ${initiative.name}`}><X size={15} /></button>
              </> : <span className="plan-slot-empty">Свободный слот</span>}
            </div>
          );
        })}
      </div>

      <div className="plan-budget-summary">
        <div><span>Использовано</span><strong>{simulation.budgetUsed} <small>/ {catalog.budget}</small></strong></div>
        <div className="plan-budget-track"><i style={{ width: `${Math.min(100, (simulation.budgetUsed / catalog.budget) * 100)}%` }} /></div>
      </div>

      {simulation.violations.length > 0 && decisions.length > 0 && (
        <div className="plan-notices">{simulation.violations.map((violation) => <p key={violation.code}><CircleAlert size={14} /> {violation.message}</p>)}</div>
      )}
      {error && <div className="plan-api-error" role="alert"><CircleAlert size={15} /><span>{error}</span><button type="button" onClick={onRetry}>Повторить</button></div>}

      <Button className="plan-submit" disabled={!simulation.submittable || explaining || simulating || Boolean(error)} onClick={onSubmit}>
        {explaining ? <><LoaderCircle className="spin" /> Готовим разбор…</> : simulating ? <><LoaderCircle className="spin" /> Считаем…</> : <>{cta}<ArrowRight size={17} /></>}
      </Button>
      <p className="plan-helper">Нужно ровно пять решений в пределах бюджета и не более двух мер одного направления.</p>
    </aside>
  );
}
