export type IndicatorCode = "T1" | "T2" | "E1" | "E2" | "S1" | "S2" | "B1" | "B2" | "C1" | "C2";
export type Direction = "transport" | "ecology" | "social" | "safety" | "services";

export interface District {
  id: string;
  name: string;
  profile: string;
  population: number;
  indicators: Record<IndicatorCode, number>;
}

export interface Initiative {
  id: string;
  direction: Direction;
  name: string;
  type: "district" | "city";
  cost: number;
  lag: number;
  effects: Partial<Record<IndicatorCode, number>>;
}

export interface CatalogResponse {
  budget: number;
  requiredDecisions: number;
  districts: District[];
  initiatives: Initiative[];
}

export interface Decision { initiativeId: string; districtId?: string }
export interface Violation { code: string; message: string }

export interface IndicatorResult {
  indicator: IndicatorCode;
  before: number;
  after: number;
  delta: number;
}

export interface DistrictResult {
  id: string;
  name: string;
  scoreBefore: number;
  scoreAfter: number;
  indicators?: IndicatorResult[];
}

export interface SimulationResponse {
  submittable: boolean;
  violations: Violation[];
  budgetUsed: number;
  budgetRemaining: number;
  districts: DistrictResult[];
  dAvgBefore: number;
  dAvgAfter: number;
  minDistrictId: string;
  minDBefore: number;
  minDAfter: number;
  criticalBefore: number;
  criticalAfter: number;
  score: number | null;
  synergies: Array<{ label: string; districtId: string }>;
  contributions: unknown[];
}

export interface HealthResponse { status: "ok" }

export interface OptimumResponse {
  bestScore: number;
  decisions?: Decision[];
}

export interface ExplainResponse {
  result: unknown;
  explanation: { summary: string; strengths: string[]; risks: string[]; recommendations: string[] };
  source: "live" | "cached";
}

export interface SavedScenario {
  id: string;
  label: string;
  decisions: Decision[];
  result: SimulationResponse;
  createdAt: string;
}

export type SimulationResult = Omit<SimulationResponse, "submittable" | "violations">;

export interface CompareScenarioResult {
  label: string;
  result?: SimulationResult;
  violations?: Violation[];
}

export interface CompareResponse {
  scenarios: CompareScenarioResult[];
  comparison: { summary: string; source: "live" | "cached" };
}
