"use client";

import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { Activity, ArrowLeft, BookmarkPlus, ClipboardList, LoaderCircle, MapPin, PanelsTopLeft, TrendingDown, TrendingUp, X } from "lucide-react";

import { AstanaThreeMap } from "@/components/astana-three-map";
import { ScenarioCompare } from "@/components/scenario-compare";
import { Button } from "@/components/ui/button";
import { api, USE_MOCKS } from "@/lib/api-client";
import type { Decision, Direction, District, DistrictResult, ExplainResponse, Initiative, OptimumResponse, SavedScenario, SimulationResponse, CatalogResponse } from "@/lib/types";
import { DistrictPicker } from "@/components/simulator/district-picker";
import { InitiativeCatalog } from "@/components/simulator/initiative-catalog";
import { PlanTray } from "@/components/simulator/plan-tray";
import { ScenarioHeader } from "@/components/simulator/scenario-header";

type DirectionFilter = "all" | Direction;
type MobileView = "catalog" | "district" | "plan";
const SCENARIO_STORAGE_KEY = "bessheshim-scenarios-v1";

export function Simulator() {
  const [catalog, setCatalog] = useState<CatalogResponse | null>(null);
  const [health, setHealth] = useState<"loading" | "ok" | "error">("loading");
  const [decisions, setDecisions] = useState<Decision[]>([]);
  const [simulation, setSimulation] = useState<SimulationResponse | null>(null);
  const [loading, setLoading] = useState(true);
  const [simulating, setSimulating] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [explanation, setExplanation] = useState<ExplainResponse | null>(null);
  const [optimum, setOptimum] = useState<OptimumResponse | null>(null);
  const [explaining, setExplaining] = useState(false);
  const [retryNonce, setRetryNonce] = useState(0);
  const [query, setQuery] = useState("");
  const [activeDirection, setActiveDirection] = useState<DirectionFilter>("all");
  const [pendingInitiative, setPendingInitiative] = useState<Initiative | null>(null);
  const [mobileView, setMobileView] = useState<MobileView>("catalog");
  const [selectedDistrictId, setSelectedDistrictId] = useState<string | null>(null);
  const [showHelp, setShowHelp] = useState(false);
  const [savedScenarios, setSavedScenarios] = useState<SavedScenario[]>([]);
  const [scenarioLabel, setScenarioLabel] = useState("Сценарий 1");
  const [saveMessage, setSaveMessage] = useState<string | null>(null);
  const [showCompare, setShowCompare] = useState(false);
  const decisionVersion = useRef(0);

  const load = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const [nextCatalog, nextHealth] = await Promise.all([api.catalog(), api.health()]);
      setCatalog(nextCatalog);
      setHealth(nextHealth.status === "ok" ? "ok" : "error");
      setSimulation(await api.simulate([]));
      setSelectedDistrictId((current) => current && nextCatalog.districts.some((district) => district.id === current) ? current : nextCatalog.districts.find((district) => district.id === "nura")?.id ?? nextCatalog.districts[0]?.id ?? null);
    } catch (reason) {
      setHealth("error");
      setError(reason instanceof Error ? reason.message : "Не удалось загрузить симулятор");
    } finally {
      setLoading(false);
    }
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
      } finally {
        if (active) setSimulating(false);
      }
    }, 120);
    return () => { active = false; window.clearTimeout(timer); };
  }, [catalog, decisions, retryNonce]);

  const selectedIds = useMemo(() => new Set(decisions.map((item) => item.initiativeId)), [decisions]);
  const directionCounts = useMemo(() => {
    if (!catalog) return {};
    return decisions.reduce<Partial<Record<Direction, number>>>((counts, decision) => {
      const direction = catalog.initiatives.find((initiative) => initiative.id === decision.initiativeId)?.direction;
      if (direction) counts[direction] = (counts[direction] ?? 0) + 1;
      return counts;
    }, {});
  }, [catalog, decisions]);

  function addDecision(initiativeId: string, districtId?: string) {
    if (selectedIds.has(initiativeId) || !catalog || decisions.length >= catalog.requiredDecisions) return;
    decisionVersion.current += 1;
    setSimulating(true);
    setExplanation(null);
    setOptimum(null);
    setError(null);
    setDecisions((current) => [...current, { initiativeId, districtId }]);
    setPendingInitiative(null);
  }

  function requestAdd(initiative: Initiative) {
    if (initiative.type === "district" && selectedDistrictId) addDecision(initiative.id, selectedDistrictId);
    else if (initiative.type === "district") setPendingInitiative(initiative);
    else addDecision(initiative.id);
  }

  function removeDecision(initiativeId: string) {
    decisionVersion.current += 1;
    setSimulating(true);
    setExplanation(null);
    setOptimum(null);
    setDecisions((current) => current.filter((item) => item.initiativeId !== initiativeId));
  }

  function reset() {
    decisionVersion.current += 1;
    setSimulating(true);
    setDecisions([]);
    setExplanation(null);
    setOptimum(null);
    setError(null);
    setMobileView("catalog");
  }

  async function submit() {
    if (!simulation?.submittable || simulating || error) return;
    const version = decisionVersion.current;
    setExplaining(true);
    setError(null);
    try {
      const [nextExplanation, nextOptimum] = await Promise.all([api.explain(decisions), api.optimum()]);
      if (version === decisionVersion.current) {
        setExplanation(nextExplanation);
        setOptimum(nextOptimum);
      }
    } catch (reason) {
      if (version === decisionVersion.current) setError(reason instanceof Error ? reason.message : "Не удалось получить разбор");
    } finally {
      setExplaining(false);
    }
  }

  function persistScenarios(next: SavedScenario[]) {
    setSavedScenarios(next);
    window.localStorage.setItem(SCENARIO_STORAGE_KEY, JSON.stringify(next));
  }

  function saveScenario() {
    const label = scenarioLabel.trim();
    if (!simulation?.submittable || !simulation.score) return;
    if (!label) { setSaveMessage("Добавьте название сценария"); return; }
    if (savedScenarios.length >= 3) { setSaveMessage("Можно сохранить не более трёх сценариев"); return; }
    const next = [...savedScenarios, {
      id: window.crypto.randomUUID(), label,
      decisions: structuredClone(decisions), result: structuredClone(simulation),
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

  if (loading) return <div className="cockpit-loading"><span className="cockpit-brand-mark">BS/5</span><LoaderCircle className="spin" /> Загружаем городскую модель…</div>;
  if (!catalog || !simulation) return <div className="cockpit-loading"><p>{error ?? "Каталог пока пуст"}</p><Button onClick={() => void load()}>Повторить</Button></div>;

  const selectedDistrict = catalog.districts.find((district) => district.id === selectedDistrictId);
  const selectedResult = simulation.districts.find((district) => district.id === selectedDistrictId);

  return (
    <div className="play-map-page" id="simulator">
      <AstanaThreeMap districts={simulation.districts} selectedDistrictId={selectedDistrictId} onSelectDistrict={(id) => { setSelectedDistrictId(id); setMobileView("district"); }} />
      <ScenarioHeader
        budget={catalog.budget}
        budgetRemaining={simulation.budgetRemaining}
        decisionCount={decisions.length}
        requiredDecisions={catalog.requiredDecisions}
        health={health}
        isMock={USE_MOCKS}
        onReset={reset}
        onHelp={() => setShowHelp(true)}
        savedCount={savedScenarios.length}
        onCompare={() => setShowCompare(true)}
      />

      <main className="play-map-layout" id="cockpit-guide">
        <aside className={`play-map-sidebar play-map-sidebar-left ${mobileView === "catalog" ? "is-mobile-active" : ""}`}>
          <div className="play-map-sidebar-caption"><span className="cockpit-kicker">АСТАНА / ГОРОДСКАЯ ЛАБОРАТОРИЯ</span><strong>Выберите меры для города</strong><p>Выберите район на карте и вложите бюджет в его развитие.</p></div>
          <div className="play-map-catalog-scroll">
            <InitiativeCatalog
              activeDirection={activeDirection}
              budgetRemaining={simulation.budgetRemaining}
              directionCounts={directionCounts}
              initiatives={catalog.initiatives}
              query={query}
              requiredDecisions={catalog.requiredDecisions}
              selectedIds={selectedIds}
              selectedCount={decisions.length}
              onDirectionChange={setActiveDirection}
              onQueryChange={setQuery}
              onRequestAdd={requestAdd}
              onRemove={removeDecision}
            />
          </div>
        </aside>

        <aside className="play-map-sidebar play-map-sidebar-right">
          <div className={`play-map-mobile-panel ${mobileView === "district" ? "is-mobile-active" : ""}`}>
            <DistrictDetails district={selectedDistrict} result={selectedResult} decisions={decisions} initiatives={catalog.initiatives} cityScore={simulation.dAvgAfter} simulating={simulating} />
          </div>
          <div className={`play-map-mobile-panel ${mobileView === "plan" ? "is-mobile-active" : ""}`}>
            <PlanTray catalog={catalog} decisions={decisions} error={error} explaining={explaining} simulating={simulating} simulation={simulation} onRemove={removeDecision} onReset={reset} onRetry={() => setRetryNonce((value) => value + 1)} onSubmit={() => void submit()} />
          </div>
        </aside>

        <div className="play-map-hint"><MapPin size={15} /><span>Перетащите карту, нажмите на район для деталей</span></div>

        <nav className="play-map-mobile-tabs" aria-label="Разделы симулятора">
          <button className={mobileView === "catalog" ? "is-active" : ""} type="button" onClick={() => setMobileView("catalog")}><PanelsTopLeft size={16} /> Меры</button>
          <button className={mobileView === "district" ? "is-active" : ""} type="button" onClick={() => setMobileView("district")}><MapPin size={16} /> Район</button>
          <button className={mobileView === "plan" ? "is-active" : ""} type="button" onClick={() => setMobileView("plan")}><ClipboardList size={16} /> План <span>{decisions.length}</span></button>
        </nav>
      </main>

      {pendingInitiative && <DistrictPicker districts={catalog.districts} initiative={pendingInitiative} onCancel={() => setPendingInitiative(null)} onSelect={(districtId) => addDecision(pendingInitiative.id, districtId)} />}
      {explanation && <ScenarioResultOverlay explanation={explanation} optimum={optimum} simulation={simulation} scenarioLabel={scenarioLabel} saveMessage={saveMessage} canSave={savedScenarios.length < 3} onLabelChange={(value) => { setScenarioLabel(value); setSaveMessage(null); }} onSave={saveScenario} onCompare={() => { setExplanation(null); setShowCompare(true); }} onClose={() => setExplanation(null)} />}
      {showHelp && <div className="play-map-help-backdrop" role="presentation" onMouseDown={(event) => { if (event.currentTarget === event.target) setShowHelp(false); }}><section className="play-map-help" role="dialog" aria-modal="true" aria-labelledby="play-map-help-title"><button type="button" onClick={() => setShowHelp(false)} aria-label="Закрыть помощь"><X size={18} /></button><span className="cockpit-kicker">КАК ИГРАТЬ</span><h2 id="play-map-help-title">Пять решений. Один город.</h2><ol><li>Выберите район на карте, чтобы увидеть его состояние.</li><li>Добавьте инициативы из панели слева. Районная мера применяется к выбранному району.</li><li>Следите за бюджетом и планом справа. После пяти решений откройте результат.</li></ol></section></div>}
      {showCompare && <div className="play-map-compare-backdrop" role="presentation" onMouseDown={(event) => { if (event.currentTarget === event.target) setShowCompare(false); }}><section className="play-map-compare-panel" role="dialog" aria-modal="true" aria-labelledby="compare-title"><div className="play-map-compare-top"><span className="cockpit-kicker">СОХРАНЁННЫЕ СЦЕНАРИИ</span><button type="button" onClick={() => setShowCompare(false)} aria-label="Закрыть сравнение"><X size={19} /></button></div><div className="compare-section" id="compare"><div className="compare-heading"><div className="eyebrow dark">A/B/C / ЛОКАЛЬНЫЕ СЦЕНАРИИ</div><h2 id="compare-title">Сравните решения</h2></div><ScenarioCompare scenarios={savedScenarios} onRemove={removeScenario} /></div></section></div>}
    </div>
  );
}

function DistrictDetails({ district, result, decisions, initiatives, cityScore, simulating }: { district?: District; result?: DistrictResult; decisions: Decision[]; initiatives: Initiative[]; cityScore: number; simulating: boolean }) {
  if (!district || !result) return <section className="play-map-district-card"><p>Выберите район на карте, чтобы увидеть его показатели.</p></section>;
  const delta = result.scoreAfter - result.scoreBefore;
  const measures = decisions.filter((decision) => decision.districtId === district.id).map((decision) => initiatives.find((initiative) => initiative.id === decision.initiativeId)).filter((initiative): initiative is Initiative => Boolean(initiative));
  return (
    <section className="play-map-district-card" aria-labelledby="play-map-district-title">
      <div className="play-map-district-top"><span className="cockpit-kicker"><MapPin size={12} /> ВЫБРАННЫЙ РАЙОН</span>{simulating ? <LoaderCircle className="spin" size={16} /> : <Activity size={16} />}</div>
      <h2 id="play-map-district-title">{district.name}</h2>
      <p className="play-map-district-profile">{district.profile}</p>
      <div className="play-map-district-score"><div><span>Индекс района</span><strong>{result.scoreAfter.toFixed(1)}</strong></div><span className={delta >= 0 ? "is-positive" : "is-negative"}>{delta >= 0 ? <TrendingUp size={15} /> : <TrendingDown size={15} />}{delta >= 0 ? "+" : ""}{delta.toFixed(1)}</span></div>
      <div className="play-map-district-comparison"><span>Средний по городу</span><strong>{cityScore.toFixed(1)}</strong></div>
      <div className="play-map-district-measures"><span>Меры в районе <strong>{measures.length}</strong></span>{measures.length ? measures.map((measure) => <p key={measure.id}>{measure.name}</p>) : <p>Пока нет мер. Выберите инициативу слева, чтобы развить этот район.</p>}</div>
    </section>
  );
}

function ScenarioResultOverlay({ explanation, optimum, simulation, scenarioLabel, saveMessage, canSave, onLabelChange, onSave, onCompare, onClose }: { explanation: ExplainResponse; optimum: OptimumResponse | null; simulation: SimulationResponse; scenarioLabel: string; saveMessage: string | null; canSave: boolean; onLabelChange: (value: string) => void; onSave: () => void; onCompare: () => void; onClose: () => void }) {
  const benchmark = optimum && simulation.score ? Math.min(100, (simulation.score / optimum.bestScore) * 100) : null;
  return <div className="play-map-result-backdrop" role="presentation" onMouseDown={(event) => { if (event.currentTarget === event.target) onClose(); }}>
    <section className="play-map-result" role="dialog" aria-modal="true" aria-labelledby="play-map-result-title">
      <div className="play-map-result-head"><span className="cockpit-kicker">РАЗБОР СЦЕНАРИЯ · {explanation.source === "cached" ? "КЭШИРОВАННЫЙ" : "LIVE AI"}</span><button type="button" onClick={onClose} aria-label="Закрыть результат"><X size={19} /></button></div>
      <h2 id="play-map-result-title">Город ответил.</h2>
      <div className="play-map-result-score"><span>Итоговый индекс</span><strong>{simulation.score?.toFixed(2)}</strong>{benchmark !== null && <small>{benchmark.toFixed(1)}% от лучшего сценария</small>}</div>
      <p className="play-map-result-summary">{explanation.explanation.summary}</p>
      <ResultList title="Что сработало" items={explanation.explanation.strengths} />
      <ResultList title="Риски" items={explanation.explanation.risks} />
      <ResultList title="Следующий шаг" items={explanation.explanation.recommendations} />
      <div className="save-scenario"><div><span>Сохранить для сравнения</span><p>Сценарий останется только в этом браузере.</p></div><label><span className="sr-only">Название сценария</span><input value={scenarioLabel} onChange={(event) => onLabelChange(event.target.value)} maxLength={80} /></label><Button onClick={onSave} disabled={!canSave}><BookmarkPlus size={16} /> Сохранить</Button>{saveMessage && <small>{saveMessage}</small>}</div>
      {saveMessage === "Сценарий сохранён" && <button className="play-map-result-back" type="button" onClick={onCompare}>Открыть сравнение</button>}
      <button className="play-map-result-back" type="button" onClick={onClose}><ArrowLeft size={16} /> Вернуться к карте</button>
    </section>
  </div>;
}

function ResultList({ title, items }: { title: string; items: string[] }) {
  return <div className="play-map-result-list"><h3>{title}</h3>{items.map((item, index) => <p key={`${title}-${index}`}>{item}</p>)}</div>;
}
