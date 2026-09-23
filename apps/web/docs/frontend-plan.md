# Frontend plan — «Аким на 5 часов»

Status: planning doc only, nothing here is implemented yet. Owner: Frontend (`apps/web/**`). Companion docs: `.ai/PLAN.md` (product), `contracts/api.md` (endpoints), `contracts/scoring.md` (formula), `.ai/DECISIONS.md`, `.ai/WORKFLOW.md`.

## Where things stand

The vertical slice on `fe/live-simulator` is uncommitted (`apps/` is untracked in `git status`). It has one route (`/`) with catalog → 5 picks → live SVG map + score → `/explain`, all inside a single 167-line `components/simulator.tsx`. Known gaps against the docs:

- README promises a live Score before the set of 5 is valid; the contract returns `score: null` until `submittable`, and the current UI substitutes `dAvgAfter` instead — a different number under the same label.
- The mock `simulate()` in `lib/api-client.ts` reimplements the scoring math but only partly (no synergies, no incompatibility/district checks, empty `contributions`), which drifts from DECISIONS #11 ("frontend never reimplements the math").
- The map colors districts by score but doesn't flag critical (<40) indicators, which is the specified killer feature.
- The result headline is hardcoded text, not derived from the outcome.
- `request()` has one 8s timeout for every call, but `/explain` and `/compare` are documented to take up to `LLM_TIMEOUT_SECONDS` (default 25s) plus retry.
- Everything lives on `/` — there's no result URL to share/reload, no compare view, no methodology page.

The backend isn't live yet (pre-Checkpoint 2), so the plan below is mock-first throughout.

## Decisions

1. **Live Score before submit.** The FE derives a *provisional* Score client-side: `0.7·dAvgAfter + 0.3·minDAfter − 1.0·criticalAfter`, using the same coefficients as `contracts/scoring.md` (sourced from `data/rules.json`, not hardcoded). It's labeled «предварительно» in the UI until the set is `submittable`, at which point the server's real `score` takes over. This reuses only the final formula line already documented, not the engine — engine logic still lives in `internal/scoring` only.
2. **Mocks are fixture-only.** No TS reimplementation of scoring. Mock `simulate()` returns real fixtures (`contracts/examples/simulate-partial.json` for 1–4 picks, `simulate-valid.json` for 5, a derived before=after baseline for 0 picks), so numbers for the fixed demo scenario are exact and nothing can drift from the Go engine. Numbers for arbitrary off-script picks in mock mode will be approximate — acceptable since mock mode is a dev/offline aid, not the judged path.
3. **District assignment via map.** Click a district-type measure to "arm" it, then click a district on the map to place it. An inline district-chip row is the fallback for keyboard use, mobile, and when the map isn't visible. Ties the required interaction to the killer feature instead of a plain `<select>`.
4. **Map stays SVG.** OSM currently returns 6 Astana districts against the contract's 5 (documented in `.ai/handoffs/frontend.md`), so this repeats DECISIONS #10's fallback. Scope is: critical-count badge per district, a delta chip on each pick, a synergy indicator on the anchored district, and a drill-down (10 indicators, before→after bars, 40-line) on click.
5. **Pages:** `/`, `/play`, `/result`, `/compare`, `/about`. Rationale below.
6. **Scenario state lives in the URL** on `/play` and `/result`: `?d=M7:nura,M8:nura,M10:nura,M12,M5:saryarka&e=<eventId>`. Survives refresh, is shareable, and lets the demo pre-bake a link to the brief's own example. Saved A/B/C comparison scenarios live in `localStorage` (browser-only per DECISIONS #12, no backend persistence).
7. **Extras build order:** optimum gauge → compare → crisis event — the reverse of the cut order in `.ai/DECISIONS.md`/`WORKFLOW.md`, so the riskiest, most contract-dependent piece (the event) is attempted last and is the first to be dropped if time runs out.
8. **`/explain` loading state is an honest two-phase UI**, not a scripted fake step list (there's genuinely only one request in flight): phase 1 shows the Score reveal instantly using the `/simulate` result already in hand; phase 2 is a skeleton labeled "AI анализирует… Ns" with an elapsed-time counter, capped at the documented 25s timeout plus one retry.
9. **ObsidianUI** components (installed via `npx shadcn add`, brings in the `motion` dependency), in priority order: Score ticker → result-card reveal → landing hero text reveal → landing hero background. The background is the first cut if time is short. All four must respect `prefers-reduced-motion`.

## Page map

| Route | Purpose | Reads | Key pieces |
|---|---|---|---|
| `/` — Landing | 10-second pitch | static copy + `/health` for a status dot | Hero (ObsidianUI text + background), 3-step explainer (бюджет → 5 решений → AI-разбор), baseline teaser ("Город сегодня: 52.56, Нура в красной зоне"), CTA → `/play`, secondary "Пример из ТЗ" deep-link → `/play?d=…` |
| `/play` — Simulator | The killer feature | `/catalog`, `/simulate` (debounced per change), `/optimum` (score only), optionally `/events/draw` | Catalog panel (arm→place), SVG map + drill-down, live panel (provisional Score, budget, critical count, violations, optimum %), event banner if shipped. Submit routes to `/result?d=…` |
| `/result` — Result | Final Score + AI breakdown | `/simulate` (instant phase 1) then `/explain` (phase 2); `/optimum?reveal=true` only on click | Score reveal, delta vs. baseline, synergies fired, AI summary/strengths/risks/recommendations, `cached` badge, optimum %, "Показать лучший набор", "Сохранить как A/B/C", "Изменить решения" → back to `/play?d=same`. Invalid/missing `d` shows violations + link back |
| `/compare` — Compare | 2–3 saved scenarios | `localStorage` + `POST /compare` | Per-scenario column (Score, budget, critical count, mini district bars), open/delete, AI comparison summary + source badge. <2 saved → empty state |
| `/about` — Methodology | Score traceability (judging criteria: value, tech) | `/catalog` + rules | Formula steps from `contracts/scoring.md`, 10-indicator weight table, district baselines (<40 flagged), 14-initiative table (cost/lag/realized-%/effects), synergies & incompatibilities, explicit "AI never computes numbers" note |

Shared layout: wordmark nav (Симулятор / Сравнение (n) / Методика), API status dot, a visible "MOCK" chip when `USE_MOCKS` is on, per-route `loading.tsx`/`error.tsx` with retry, and `not-found.tsx`. `/play` and `/result` read `useSearchParams`, which needs a `<Suspense>` boundary in this Next version — check `node_modules/next/dist/docs/` before writing route code (per `apps/web/AGENTS.md`, this Next release has breaking changes from training-data defaults).

Demo path: `/` → "Пример из ТЗ" → submit → `/result` (56.54, +3.99, AI, optimum %) → save A → tweak → save B → `/compare`.

## Build sequence (for whoever implements this next)

0. Commit the existing slice (`apps/web`, not `.next/`) before anything else — it's currently untracked and at risk.
1. **Foundation:** fix `lib/types.ts` shapes (`synergies`, `contributions`, `ExplainResponse.result`, add `OptimumResponse`/`Compare*`/`CityEvent`/`ApiErrorBody`); rework `lib/api-client.ts` (per-call timeouts, `optimum`/`compare`/`drawEvent`, fixture-only mocks); add `lib/score.ts` (provisional score) and `lib/scenario-url.ts` (encode/decode `?d=`); extract `hooks/use-scenario.ts` from `simulator.tsx`.
2. **Routing + shared layout:** turn `/` into the landing page, move the simulator to `/play`, stub `/result`, `/compare`, `/about` with their own loading/error/empty states, add `site-nav.tsx`.
3. **`/play` interaction:** split `catalog-panel.tsx`, extend `city-map.tsx` (arm/place, keyboard support, critical badges, deltas, synergy marker), add `district-drilldown.tsx`, add `live-panel.tsx`.
4. **`/result`:** `result-panel.tsx` with the two-phase reveal described in decision 8, error/retry, 422 handling.
5. **Optimum gauge:** live on `/play`, official on `/result`, reveal button behind submit.
6. **`/compare` + `/about`:** `lib/scenarios.ts` (localStorage, max 3, try/catch everywhere), `compare-view.tsx`, static `/about`.
7. **ObsidianUI polish**, then responsive QA at 1280px/390px with reduced motion.
8. **Crisis event** on `/play`, only if `eventId` support lands on the backend contract (see below) by Checkpoint 3; otherwise cut and note it in the handoff.

## Coordination needed (to post under `NEEDS:` in `.ai/handoffs/frontend.md`)

- **BE + AI+data:** an optional `eventId` on `POST /simulate`, `/explain`, `/compare` so the server (not the client) applies `shocks`/`budgetDelta` before scoring — today's contract has no way to make the event affect the numbers consistently across map/score/AI.
- **BE + AI+data (nice-to-have):** `/catalog` also returning `rules` (weights, score coefficients, synergies, incompatibilities, critical threshold) so `/about` and the provisional-score formula read only from the API, matching the README's "читается фронтом только через API." Until then, `/about` and `lib/score.ts` read `data/rules.json` directly at build time.
- **Whoever owns the README:** note that the live Score is provisional until 5 valid picks, and add the new page list + demo deep-link to the "Сценарий демо" section.

## Verification checklist (for implementation time, not now)

- `npm run lint && npm run build` after each step.
- Mock-mode walk-through of the full demo path above, confirming 52.56 baseline, 56.54 final, cached badge, optimum %, and that `/result` survives a refresh.
- Bad input: 6th pick disabled, budget-blocked cards disabled, garbage `?d=` on `/result` shows violations, unplaced armed measure blocks submit, M1+M3 incompatibility surfaces once the real API is up.
- Error paths: API down → retry on every route; slow `/explain` → elapsed counter then timeout+retry; unknown route → 404.
- 1280px and 390px, `prefers-reduced-motion` on, keyboard-only district assignment.
- Once BE is live: `NEXT_PUBLIC_USE_MOCKS=false`, repeat the walk-through, cross-check against `data/fixtures/golden.json`.
