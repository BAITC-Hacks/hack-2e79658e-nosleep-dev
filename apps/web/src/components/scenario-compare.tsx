"use client";

import { useMemo, useState } from "react";
import { Check, CircleAlert, LoaderCircle, Trash2 } from "lucide-react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { api } from "@/lib/api-client";
import type { CompareResponse, SavedScenario } from "@/lib/types";

export function ScenarioCompare({ scenarios, onRemove }: { scenarios: SavedScenario[]; onRemove: (id: string) => void }) {
  const [selected, setSelected] = useState<string[]>(scenarios.map((scenario) => scenario.id).slice(0, 3));
  const [comparison, setComparison] = useState<CompareResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const selectedScenarios = useMemo(
    () => selected.flatMap((id) => scenarios.find((scenario) => scenario.id === id) ?? []),
    [scenarios, selected],
  );

  function toggle(id: string) {
    setComparison(null);
    setSelected((current) => current.includes(id) ? current.filter((item) => item !== id) : current.length < 3 ? [...current, id] : current);
  }

  async function compare() {
    if (selectedScenarios.length < 2) return;
    setLoading(true);
    setError(null);
    try {
      setComparison(await api.compare(selectedScenarios.map(({ label, decisions }) => ({ label, decisions }))));
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось сравнить сценарии");
    } finally {
      setLoading(false);
    }
  }

  if (scenarios.length === 0) {
    return <div className="compare-empty">Сохраните первый сценарий — здесь появится сравнение.</div>;
  }

  return (
    <>
      <div className="saved-scenarios">
        {scenarios.map((scenario) => {
          const active = selected.includes(scenario.id);
          return (
            <article className={`saved-scenario ${active ? "active" : ""}`} key={scenario.id}>
              <button className="scenario-select" onClick={() => toggle(scenario.id)} aria-pressed={active}>
                <span className="scenario-check">{active && <Check size={13} />}</span>
                <span><small>Сценарий</small><strong>{scenario.label}</strong></span>
              </button>
              <div className="scenario-metrics">
                <span><small>Score</small><strong>{scenario.result.score?.toFixed(2) ?? "—"}</strong></span>
                <span><small>Бюджет</small><strong>{scenario.result.budgetUsed}</strong></span>
                <span><small>Критично</small><strong>{scenario.result.criticalAfter}</strong></span>
              </div>
              <button className="scenario-delete" onClick={() => onRemove(scenario.id)} aria-label={`Удалить ${scenario.label}`}><Trash2 size={15} /></button>
            </article>
          );
        })}
      </div>

      <div className="compare-actions">
        <p>Выберите 2–3 сценария. Данные остаются только в этом браузере.</p>
        <Button onClick={() => void compare()} disabled={selectedScenarios.length < 2 || loading}>
          {loading ? <><LoaderCircle className="spin" /> Сравниваем…</> : `Сравнить ${selectedScenarios.length || ""}`}
        </Button>
      </div>

      {error && <div className="compare-error"><CircleAlert size={16} />{error}<button onClick={() => void compare()}>Повторить</button></div>}

      {comparison && (
        <div className="comparison-result">
          <div className="comparison-summary">
            <div><span>AI / СРАВНЕНИЕ</span><Badge>{comparison.comparison.source === "cached" ? "Кэшировано" : "Live AI"}</Badge></div>
            <p>{comparison.comparison.summary}</p>
          </div>
          <div className="comparison-table">
            {comparison.scenarios.map((scenario) => (
              <article key={scenario.label}>
                <h3>{scenario.label}</h3>
                {scenario.result ? (
                  <>
                    <strong>{scenario.result.score?.toFixed(2) ?? "—"}</strong>
                    <dl><div><dt>Бюджет</dt><dd>{scenario.result.budgetUsed}/100</dd></div><div><dt>Средний балл</dt><dd>{scenario.result.dAvgAfter.toFixed(2)}</dd></div><div><dt>Критические</dt><dd>{scenario.result.criticalAfter}</dd></div></dl>
                  </>
                ) : <div className="scenario-invalid"><CircleAlert size={15} />{scenario.violations?.[0]?.message ?? "Сценарий недействителен"}</div>}
              </article>
            ))}
          </div>
        </div>
      )}
    </>
  );
}
