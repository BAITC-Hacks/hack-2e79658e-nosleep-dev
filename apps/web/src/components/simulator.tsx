"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { ArrowRight, BookmarkPlus, Check, ChevronRight, CircleAlert, LoaderCircle, RotateCcw, X } from "lucide-react";

import { CityMap } from "@/components/city-map";
import { ScenarioCompare } from "@/components/scenario-compare";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { api, USE_MOCKS } from "@/lib/api-client";
import type { CatalogResponse, Decision, ExplainResponse, OptimumResponse, SavedScenario, SimulationResponse } from "@/lib/types";

const directionLabels = { transport: "Транспорт", ecology: "Экология", social: "Соцсфера", safety: "Безопасность", services: "Сервисы" };
const SCENARIO_STORAGE_KEY = "akim-5h-scenarios-v1";

export function Simulator() {
  const [catalog, setCatalog] = useState<CatalogResponse | null>(null);
  const [health, setHealth] = useState<"loading" | "ok" | "error">("loading");
  const [decisions, setDecisions] = useState<Decision[]>([]);
  const [districtDrafts, setDistrictDrafts] = useState<Record<string, string>>({});
  const [simulation, setSimulation] = useState<SimulationResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [simulating, setSimulating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [explanation, setExplanation] = useState<ExplainResponse | null>(null);
  const [optimum, setOptimum] = useState<OptimumResponse | null>(null);
  const [explaining, setExplaining] = useState(false);
  const [retryNonce, setRetryNonce] = useState(0);
  const [savedScenarios, setSavedScenarios] = useState<SavedScenario[]>([]);
  const [scenarioLabel, setScenarioLabel] = useState("Сценарий 1");
  const [saveMessage, setSaveMessage] = useState<string | null>(null);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [nextCatalog, nextHealth] = await Promise.all([api.catalog(), api.health()]);
      setCatalog(nextCatalog);
      setHealth(nextHealth.status === "ok" ? "ok" : "error");
      setSimulation(await api.simulate([]));
    } catch (reason) {
      setHealth("error");
      setError(reason instanceof Error ? reason.message : "Не удалось загрузить симулятор");
    } finally { setLoading(false); }
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => void load(), 0);
    return () => window.clearTimeout(timer);
  }, [load]);

  useEffect(() => {
    const timer = window.setTimeout(() => {
      try {
        const stored = JSON.parse(window.localStorage.getItem(SCENARIO_STORAGE_KEY) ?? "[]");
        if (Array.isArray(stored)) {
          const next = stored.slice(0, 3) as SavedScenario[];
          setSavedScenarios(next);
          setScenarioLabel(`Сценарий ${next.length + 1}`);
        }
      } catch {
        window.localStorage.removeItem(SCENARIO_STORAGE_KEY);
      }
    }, 0);
    return () => window.clearTimeout(timer);
  }, []);

  useEffect(() => {
    if (!catalog) return;
    let active = true;
    const timer = window.setTimeout(async () => {
      setSimulating(true);
      setError(null);
      try {
        const result = await api.simulate(decisions);
        if (active) setSimulation(result);
      } catch (reason) {
        if (active) setError(reason instanceof Error ? reason.message : "Расчёт не выполнен");
      } finally { if (active) setSimulating(false); }
    }, 120);
    return () => { active = false; window.clearTimeout(timer); };
  }, [catalog, decisions, retryNonce]);

  const selectedIds = useMemo(() => new Set(decisions.map((item) => item.initiativeId)), [decisions]);

  function addDecision(initiativeId: string, type: "district" | "city") {
    if (selectedIds.has(initiativeId) || decisions.length >= 5) return;
    const districtId = type === "district" ? districtDrafts[initiativeId] : undefined;
    if (type === "district" && !districtId) { setError("Сначала выберите район для этой меры"); return; }
    setExplanation(null);
    setOptimum(null);
    setDecisions((current) => [...current, { initiativeId, districtId }]);
  }

  function removeDecision(initiativeId: string) {
    setExplanation(null);
    setOptimum(null);
    setDecisions((current) => current.filter((item) => item.initiativeId !== initiativeId));
  }

  async function submit() {
    if (!simulation?.submittable) return;
    setExplaining(true);
    setError(null);
    try {
      const [nextExplanation, nextOptimum] = await Promise.all([api.explain(decisions), api.optimum()]);
      setExplanation(nextExplanation);
      setOptimum(nextOptimum);
      window.setTimeout(() => document.getElementById("result")?.scrollIntoView({ behavior: "smooth" }), 0);
    } catch (reason) {
      setError(reason instanceof Error ? reason.message : "Не удалось получить разбор");
    } finally { setExplaining(false); }
  }

  function persistScenarios(next: SavedScenario[]) {
    setSavedScenarios(next);
    window.localStorage.setItem(SCENARIO_STORAGE_KEY, JSON.stringify(next));
  }

  function saveScenario() {
    const label = scenarioLabel.trim();
    const currentSimulation = simulation;
    if (!currentSimulation?.submittable || !currentSimulation.score) return;
    if (!label) { setSaveMessage("Добавьте название сценария"); return; }
    if (savedScenarios.length >= 3) { setSaveMessage("Можно сохранить не более трёх сценариев"); return; }
    const next = [...savedScenarios, {
      id: window.crypto.randomUUID(),
      label,
      decisions: structuredClone(decisions),
      result: structuredClone(currentSimulation),
      createdAt: new Date().toISOString(),
    }];
    persistScenarios(next);
    setScenarioLabel(`Сценарий ${next.length + 1}`);
    setSaveMessage("Сценарий сохранён");
  }

  function removeScenario(id: string) {
    persistScenarios(savedScenarios.filter((scenario) => scenario.id !== id));
    setSaveMessage(null);
  }

  if (loading) return <div className="loading-state"><LoaderCircle className="spin" /> Загружаем городские данные…</div>;
  if (!catalog || !simulation) {
    return <div className="empty-state"><p>{error ?? "Каталог пока пуст"}</p><Button onClick={() => void load()}>Повторить</Button></div>;
  }

  return (
    <>
      <nav className="nav-bar">
        <a className="wordmark" href="#top"><span />АКИМ / 5Ч</a>
        <div className="nav-meta"><span className={`status-dot ${health}`} /> API {USE_MOCKS ? "mock" : health}<a href="#simulator">Симулятор</a>{savedScenarios.length > 0 && <a href="#compare">Сравнение</a>}</div>
      </nav>

      <header className="hero" id="top">
        <div className="eyebrow">ГОРОДСКАЯ ЛАБОРАТОРИЯ / АСТАНА</div>
        <h1>Пять решений.<br /><span>Один город.</span></h1>
        <div className="hero-bottom">
          <p>Распределите 100 единиц бюджета и смотрите, как каждый выбор меняет качество жизни районов — до того, как решение принято.</p>
          <a className="hero-cta" href="#simulator">Начать сценарий <ArrowRight size={18} /></a>
        </div>
      </header>

      <main className="simulator" id="simulator">
        <section className="sim-intro">
          <div><div className="eyebrow dark">LIVE / СЦЕНАРИЙ 01</div><h2>Соберите пакет мер</h2></div>
          <div className="budget-block"><span>Бюджет</span><strong>{simulation.budgetRemaining}</strong><small> / {catalog.budget}</small><Progress value={(simulation.budgetUsed / catalog.budget) * 100} /></div>
        </section>

        <section className="workspace">
          <div className="catalog-panel">
            <div className="panel-heading"><span>Каталог / {catalog.initiatives.length}</span><span>{decisions.length} из {catalog.requiredDecisions}</span></div>
            <div className="initiative-list">
              {catalog.initiatives.map((initiative) => {
                const selected = selectedIds.has(initiative.id);
                return (
                  <article className={`initiative ${selected ? "selected" : ""}`} key={initiative.id}>
                    <div className="initiative-copy">
                      <div className="initiative-meta"><span>{initiative.id}</span><span>{directionLabels[initiative.direction]}</span><span>{initiative.cost} ед.</span></div>
                      <h3>{initiative.name}</h3>
                      {initiative.type === "district" && !selected && (
                        <label><span className="sr-only">Район</span><select value={districtDrafts[initiative.id] ?? ""} onChange={(event) => setDistrictDrafts((current) => ({ ...current, [initiative.id]: event.target.value }))}>
                          <option value="">Выберите район</option>
                          {catalog.districts.map((district) => <option key={district.id} value={district.id}>{district.name}</option>)}
                        </select></label>
                      )}
                    </div>
                    {selected ? <button className="remove-button" onClick={() => removeDecision(initiative.id)} aria-label={`Убрать ${initiative.name}`}><X size={18} /></button> : <button className="add-button" onClick={() => addDecision(initiative.id, initiative.type)} disabled={decisions.length >= 5} aria-label={`Добавить ${initiative.name}`}><ChevronRight size={20} /></button>}
                  </article>
                );
              })}
            </div>
          </div>

          <aside className="live-panel">
            <div className="live-heading"><div><span className="pulse" /> LIVE SCORE</div>{simulating && <LoaderCircle className="spin" size={16} />}</div>
            <div className="score-lockup"><strong>{(simulation.score ?? simulation.dAvgAfter).toFixed(1)}</strong><span>{simulation.score ? "официальный Score" : "средний балл"}</span><em>{simulation.dAvgAfter >= simulation.dAvgBefore ? "+" : ""}{(simulation.dAvgAfter - simulation.dAvgBefore).toFixed(2)}</em></div>
            <CityMap districts={simulation.districts} />
            <div className="critical-strip"><span>Критические показатели</span><strong>{simulation.criticalBefore} → {simulation.criticalAfter}</strong></div>
            {simulation.violations.length > 0 && decisions.length > 0 && <div className="violations">{simulation.violations.map((violation) => <p key={violation.code}><CircleAlert size={15} />{violation.message}</p>)}</div>}
            {error && <div className="api-error" role="alert"><CircleAlert size={16} /><span>{error}</span><button onClick={() => setRetryNonce((value) => value + 1)}>Повторить</button></div>}
            <Button className="submit-button" disabled={!simulation.submittable || explaining} onClick={() => void submit()}>{explaining ? <><LoaderCircle className="spin" /> Сверяем правила и готовим разбор…</> : <>Зафиксировать решения <ArrowRight /></>}</Button>
            <button className="reset-button" onClick={() => { setDecisions([]); setExplanation(null); setOptimum(null); }} disabled={decisions.length === 0}><RotateCcw size={14} /> Сбросить сценарий</button>
          </aside>
        </section>
      </main>

      {explanation && (
        <section className="result" id="result">
          <div className="result-top"><div><div className="eyebrow">AI / РАЗБОР СЦЕНАРИЯ</div><h2>Город стал<br />сильнее.</h2></div><div className="result-score"><span>Итоговый Score</span><strong>{simulation.score?.toFixed(2)}</strong><Badge>{explanation.source === "cached" ? "Кэшированный разбор" : "Live AI"}</Badge></div></div>
          {optimum && simulation.score && <div className="result-benchmark"><div><span>От лучшего сценария</span><strong>{Math.min(100, (simulation.score / optimum.bestScore) * 100).toFixed(1)}%</strong></div><Progress value={Math.min(100, (simulation.score / optimum.bestScore) * 100)} /><p>Ваш результат {simulation.score.toFixed(2)} из оптимальных {optimum.bestScore.toFixed(2)}</p></div>}
          <p className="result-summary">{explanation.explanation.summary}</p>
          <div className="result-grid"><ResultColumn title="Что сработало" items={explanation.explanation.strengths} positive /><ResultColumn title="Риски" items={explanation.explanation.risks} /><ResultColumn title="Следующий шаг" items={explanation.explanation.recommendations} /></div>
          <div className="save-scenario">
            <div><span>Сохранить для сравнения</span><p>Сценарий останется только в этом браузере.</p></div>
            <label><span className="sr-only">Название сценария</span><input value={scenarioLabel} onChange={(event) => { setScenarioLabel(event.target.value); setSaveMessage(null); }} maxLength={80} /></label>
            <Button onClick={saveScenario} disabled={savedScenarios.length >= 3}><BookmarkPlus size={16} /> Сохранить</Button>
            {saveMessage && <small>{saveMessage}</small>}
          </div>
        </section>
      )}

      {savedScenarios.length > 0 && (
        <section className="compare-section" id="compare">
          <div className="compare-heading"><div className="eyebrow dark">A/B/C / ЛОКАЛЬНЫЕ СЦЕНАРИИ</div><h2>Сравните решения</h2></div>
          <ScenarioCompare scenarios={savedScenarios} onRemove={removeScenario} />
        </section>
      )}
    </>
  );
}

function ResultColumn({ title, items, positive = false }: { title: string; items: string[]; positive?: boolean }) {
  return <div><h3>{title}</h3>{items.map((item, index) => <p key={item}>{positive ? <Check size={16} /> : <span className="result-index">0{index + 1}</span>}{item}</p>)}</div>;
}
