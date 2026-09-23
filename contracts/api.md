# API contract — `apps/api` (Go/Gin, base path `/api/v1`)

Owner: AI+data (`internal/scoring`, `internal/llm`) + BE (`internal/http`, wiring). Changes here need approval from both sides. See `contracts/go-interfaces.md` for the Go boundary between them.

CORS: allow origin `WEB_ORIGIN` (default `http://localhost:3000`), methods `GET,POST,OPTIONS`.

## Error shape

Every non-2xx response:
```json
{ "error": { "code": "BUDGET_EXCEEDED", "message": "бюджет превышен: 127 из 100", "details": {} } }
```
`details` is optional. Codes are UPPER_SNAKE_CASE and stable — the frontend matches on `code`, not `message` (which is Russian, for display).

Violation codes (used in `simulate.violations[].code` and as the top-level error `code` on `/explain`, `/compare`):
`DECISION_COUNT`, `UNKNOWN_INITIATIVE`, `DUPLICATE_INITIATIVE`, `DISTRICT_REQUIRED`, `CITY_INITIATIVE_HAS_DISTRICT`, `UNKNOWN_DISTRICT`, `BUDGET_EXCEEDED`, `DIRECTION_LIMIT_EXCEEDED`, `INCOMPATIBLE`.

## `GET /health`
`200 { "status": "ok" }`. No auth, no deps checked — just proves the process is up.

## `GET /catalog`
Returns the full static dataset so the frontend never hardcodes it.
```json
{ "budget": 100, "requiredDecisions": 5, "districts": [ ... data/districts.json ... ], "initiatives": [ ... data/initiatives.json ... ] }
```
See `contracts/examples/catalog.json`.

## `POST /simulate`
Live-preview endpoint — called on every pick, 0 to 5 decisions. **Never calls the LLM.** Pure `internal/scoring` math.

Request:
```json
{ "decisions": [ { "initiativeId": "M7", "districtId": "nura" }, { "initiativeId": "M12" } ] }
```

Response `200` (always 200 even when invalid — the UI reads `submittable` and `violations`, it doesn't branch on HTTP status for a partial/invalid set):
```json
{
  "submittable": false,
  "violations": [ { "code": "DECISION_COUNT", "message": "нужно ровно 5 решений, получено 2" } ],
  "budgetUsed": 34, "budgetRemaining": 66,
  "districts": [ { "id": "nura", "name": "Нура", "scoreBefore": 49.18, "scoreAfter": 51.02,
    "indicators": [ { "indicator": "S1", "before": 38, "after": 40, "delta": 2 }, "... 9 more ..." ] }, "... 4 more districts ..." ],
  "dAvgBefore": 56.8624, "dAvgAfter": 57.10,
  "minDistrictId": "nura", "minDBefore": 49.18, "minDAfter": 51.02,
  "criticalBefore": 2, "criticalAfter": 1,
  "score": null,
  "synergies": [],
  "contributions": [ { "initiativeId": "M7", "name": "Школа + детсад", "direction": "social", "districtId": "nura", "cost": 24, "lag": 3, "realizedFraction": 0.625, "effectsApplied": { "S1": 10 } } ]
}
```
`score` is `null` unless `submittable` is `true` (exactly 5, no violations) — the UI shows the running district/D_avg numbers live, but the official Score only appears once the set is legal. See `contracts/examples/simulate-partial.json` and `simulate-valid.json`.

## `POST /explain`
Request: `{ "decisions": [ ...exactly 5... ] }`. Server re-validates and recomputes (client-sent numbers are never trusted).

`200`:
```json
{
  "result": { "...same shape as /simulate's body when submittable, minus violations/submittable..." },
  "explanation": { "summary": "...", "strengths": ["..."], "risks": ["..."], "recommendations": ["..."] },
  "source": "live"
}
```
`source` is `"live"` (real LLM call succeeded) or `"cached"` (DEMO_MODE, or the live call failed/timed out and the server fell back). The UI must show a small badge for `"cached"` — cached output is never presented as live. `422` with a `violations` error body if the set isn't valid. See `contracts/examples/explain-cached.json`.

## `GET /optimum`
Brute-forced once at API startup over all valid 5-decision sets (≤ a few thousand after the direction/incompatibility filters). No request body.

`200`: `{ "bestScore": 60.14 }`. With `?reveal=true`, also includes `"decisions": [...]` — the UI only reveals this after the player has submitted their own scenario (never spoils it beforehand).

## `POST /compare`
Request: `{ "scenarios": [ { "label": "A", "decisions": [...5] }, { "label": "B", "decisions": [...5] } ] }` (2–3 scenarios).

`200`: `{ "scenarios": [ { "label": "A", "result": {...} }, ... ], "comparison": { "summary": "...", "source": "live" } }`. Each scenario is validated independently; an invalid one reports its `violations` in place of a `result` and the AI comparison skips it.

## `POST /events/draw`
Request: `{ "seed": 42 }` (optional — omit for random). Response: one event from `data/events.json` (authored by AI+data, not part of this docs push):
```json
{ "id": "heating-failure-almaty", "title": "Авария теплосети", "description": "...", "districtId": "almaty", "shocks": { "C1": -15 }, "budgetDelta": -10 }
```
The frontend applies `shocks` as an immediate indicator hit (shown before the next `/simulate` call) and `budgetDelta` adjusts the working budget for the rest of the session.

## DEMO_MODE

When `DEMO_MODE=true`, `/explain` and `/compare` skip the live LLM call entirely and return the cached template response, always with `"source": "cached"`. This must work with zero API keys configured, so a dead network never kills the demo. See the Environment Variables table in `README.md`.
