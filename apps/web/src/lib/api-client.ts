import catalogExample from "../../../../contracts/examples/catalog.json";
import compareExample from "../../../../contracts/examples/compare.json";
import explainExample from "../../../../contracts/examples/explain-cached.json";
import optimumExample from "../../../../contracts/examples/optimum.json";
import validExample from "../../../../contracts/examples/simulate-valid.json";
import rules from "../../../../data/rules.json";

import type { CatalogResponse, CompareResponse, Decision, DistrictResult, ExplainResponse, HealthResponse, IndicatorCode, OptimumResponse, SimulationResponse, Violation } from "@/lib/types";

const configuredApiUrl = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000").replace(/\/$/, "");
const API_URL = configuredApiUrl.endsWith("/api/v1") ? configuredApiUrl : `${configuredApiUrl}/api/v1`;
export const USE_MOCKS = process.env.NEXT_PUBLIC_USE_MOCKS !== "false";
const catalog = catalogExample as CatalogResponse;

class ApiError extends Error {
  constructor(message: string, public readonly code = "REQUEST_FAILED") { super(message); }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const controller = new AbortController();
  const timeout = window.setTimeout(() => controller.abort(), 8_000);
  try {
    const response = await fetch(`${API_URL}${path}`, {
      ...init,
      headers: { "Content-Type": "application/json", ...init?.headers },
      signal: controller.signal,
    });
    const body = await response.json();
    if (!response.ok) throw new ApiError(body?.error?.message ?? "Сервис временно недоступен", body?.error?.code);
    return body as T;
  } catch (error) {
    if (error instanceof ApiError) throw error;
    if (error instanceof DOMException && error.name === "AbortError") throw new ApiError("Сервис не ответил за 8 секунд", "TIMEOUT");
    throw new ApiError("Не удалось связаться с API. Проверьте, запущен ли backend.");
  } finally {
    window.clearTimeout(timeout);
  }
}

function validate(decisions: Decision[]): Violation[] {
  if (decisions.length !== catalog.requiredDecisions) {
    return [{ code: "DECISION_COUNT", message: `Выберите ровно 5 мер — сейчас ${decisions.length}` }];
  }
  const chosen = decisions.map((decision) => decision.initiativeId);
  const budgetUsed = catalog.initiatives.filter((initiative) => chosen.includes(initiative.id)).reduce((sum, initiative) => sum + initiative.cost, 0);
  if (budgetUsed > catalog.budget) return [{ code: "BUDGET_EXCEEDED", message: `Бюджет превышен: ${budgetUsed} из ${catalog.budget}` }];
  const directions = chosen.map((id) => catalog.initiatives.find((initiative) => initiative.id === id)?.direction);
  if (directions.some((direction) => directions.filter((candidate) => candidate === direction).length > 2)) {
    return [{ code: "DIRECTION_LIMIT_EXCEEDED", message: "В одном направлении можно выбрать не более 2 мер" }];
  }
  return [];
}

function mockSimulation(decisions: Decision[]): SimulationResponse {
  const selected = decisions.flatMap((decision) => {
    const initiative = catalog.initiatives.find((item) => item.id === decision.initiativeId);
    return initiative ? [{ decision, initiative }] : [];
  });
  const budgetUsed = selected.reduce((sum, item) => sum + item.initiative.cost, 0);
  const weights = rules.indicatorWeights as Record<IndicatorCode, number>;
  const beforeScores = new Map((validExample.response.districts as DistrictResult[]).map((district) => [district.id, district.scoreBefore]));

  const districts = catalog.districts.map((district): DistrictResult => {
    const indicators = Object.entries(district.indicators).map(([code, before]) => {
      const indicator = code as IndicatorCode;
      let delta = 0;
      selected.forEach(({ decision, initiative }) => {
        if (initiative.type === "city" || decision.districtId === district.id) {
          delta += (initiative.effects[indicator] ?? 0) * ((8 - initiative.lag) / 8);
        }
      });
      const after = Math.max(0, Math.min(100, before + delta));
      return { indicator, before, after, delta: after - before };
    });
    const scoreAfter = indicators.reduce((sum, item) => sum + item.after * weights[item.indicator], 0);
    return { id: district.id, name: district.name, scoreBefore: beforeScores.get(district.id) ?? scoreAfter, scoreAfter, indicators };
  });

  const dAvgAfter = districts.reduce((sum, district) => sum + district.scoreAfter * (catalog.districts.find((item) => item.id === district.id)?.population ?? 0), 0);
  const worst = [...districts].sort((a, b) => a.scoreAfter - b.scoreAfter)[0];
  const criticalAfter = districts.reduce((count, district) => count + (district.indicators?.filter((indicator) => indicator.after < 40).length ?? 0), 0);
  const violations = validate(decisions);
  const score = 0.7 * dAvgAfter + 0.3 * worst.scoreAfter - criticalAfter;
  return {
    submittable: violations.length === 0,
    violations,
    budgetUsed,
    budgetRemaining: catalog.budget - budgetUsed,
    districts,
    dAvgBefore: 56.8624,
    dAvgAfter,
    minDistrictId: worst.id,
    minDBefore: 49.18,
    minDAfter: worst.scoreAfter,
    criticalBefore: 2,
    criticalAfter,
    score: violations.length === 0 ? score : null,
    synergies: [],
    contributions: [],
  };
}

export const api = {
  health: async (): Promise<HealthResponse> => USE_MOCKS ? { status: "ok" } : request<HealthResponse>("/health"),
  catalog: async (): Promise<CatalogResponse> => USE_MOCKS ? structuredClone(catalog) : request<CatalogResponse>("/catalog"),
  simulate: async (decisions: Decision[]): Promise<SimulationResponse> => USE_MOCKS ? mockSimulation(decisions) : request<SimulationResponse>("/simulate", { method: "POST", body: JSON.stringify({ decisions }) }),
  explain: async (decisions: Decision[]): Promise<ExplainResponse> => USE_MOCKS ? (structuredClone(explainExample.response) as unknown as ExplainResponse) : request<ExplainResponse>("/explain", { method: "POST", body: JSON.stringify({ decisions }) }),
  optimum: async (reveal = false): Promise<OptimumResponse> => USE_MOCKS
    ? structuredClone(reveal ? optimumExample.revealed : optimumExample.default)
    : request<OptimumResponse>(`/optimum${reveal ? "?reveal=true" : ""}`),
  compare: async (scenarios: Array<{ label: string; decisions: Decision[] }>): Promise<CompareResponse> => {
    if (!USE_MOCKS) return request<CompareResponse>("/compare", { method: "POST", body: JSON.stringify({ scenarios }) });
    return {
      scenarios: scenarios.map((scenario) => {
        const result = mockSimulation(scenario.decisions);
        return result.submittable
          ? { label: scenario.label, result }
          : { label: scenario.label, violations: result.violations };
      }),
      comparison: structuredClone(compareExample.response.comparison) as CompareResponse["comparison"],
    };
  },
};

export { ApiError };
