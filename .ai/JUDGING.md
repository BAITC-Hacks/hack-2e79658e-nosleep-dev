# JUDGING — how we score against the rubric

## Official criteria (100 pts)

### Соответствие задаче и работоспособность — 25 pts
- **Product evidence:** all 6 must-haves from the brief implemented (`.ai/PLAN.md` MVP list) — shared budget, exactly-5 decisions, budget-overrun block, AI analysis, final Score, explanation of strengths/risks/tradeoffs.
- **Demo evidence:** the 60-second flow (`DEMO.md`, to be written at T+4:15) walks through all 6 live, ending on the Score + AI explanation.
- **Weakness to watch:** if `/explain` degrades to `DEMO_MODE` during the live demo, the AI-analysis point still lands (cached ≠ absent) but say so out loud rather than let a judge notice the badge and wonder.

### Техническая реализация — 25 pts
- **Product evidence:** deterministic scoring engine (`contracts/scoring.md`) with a golden-fixture test suite (`data/fixtures/golden.json`), structured-output LLM call with schema validation and a documented retry/fallback path (`contracts/api.md` §DEMO_MODE), clean FE/BE contract (`contracts/`).
- **Demo evidence:** mention the architecture in 1 sentence during the pitch ("Go engine does the math, one LLM call explains it — the AI never invents a number").
- **Weakness to watch:** no automated tests for the frontend or for `/optimum`'s brute force — acceptable for hackathon scope, but don't claim more test coverage than exists.

### README и воспроизводимость — 25 pts
- **Product evidence:** `README.md` quick start (≤5 commands from clean clone), env var table, `DEMO_MODE` explained and working with zero keys, architecture diagram, demo scenario with the brief's own worked example numbers.
- **Demo evidence:** literally follow the README once during rehearsal on a machine that hasn't touched the repo (or `git clone` into a temp dir) before the freeze.
- **Weakness to watch:** this is the highest-leverage criterion relative to effort — don't let it be the thing rushed at T+4:50. It has its own timeline slot (T+4:15–4:45) on purpose.

### Ценность и применимость решения — 15 pts
- **Product evidence:** the scoring model is directly traceable to real urban KPIs (traffic, air quality, clinics, crime, utility reliability) with real weights, not an arbitrary game score.
- **Demo evidence:** name-check 1–2 real tradeoffs the AI surfaces (e.g. "fixing Nura's schools costs Yesil nothing, but leaves Nura's transport untouched") — shows the tool would say something a human planner would find useful.

### Потенциал развития и оригинальность подхода — 10 pts
- **Product evidence:** `/optimum` (brute-force best-possible comparison), `/compare` (multi-scenario), `/events/draw` (crisis events) are all optional-but-implemented items from the brief, not just the must-haves.
- **Demo evidence:** if time allows in the 60s, show the "% of optimum" gauge — it's the single cheapest "wow, it thinks ahead" beat.
- **Weakness to watch:** don't oversell "agentic AI" — there's one structured LLM call, not a multi-agent system. Describe it accurately.

## Brief's own "Критерии проверки" → what proves it
| Criterion | Evidence |
|---|---|
| All teams start with the same budget and data | `data/*.json` is static and versioned; no per-session randomization of inputs |
| System never allows exceeding budget | `BUDGET_EXCEEDED` violation, enforced server-side in `internal/scoring.Validate`, covered by `data/fixtures/golden.json:budget_exceeded` |
| Decisions affect the final indicators | `contracts/scoring.md` steps 2–5; visible live in `/simulate`'s before/after per indicator |
| AI gives an understandable explanation of the result and the main tradeoffs | `/explain` response's `strengths`/`risks`/`recommendations`, always grounded in the `result` it's given, never computed by the LLM |
| Changing the decision set changes the Score | Trivially true of the formula, and explicitly demonstrated by the two worked examples in `contracts/scoring.md` (base 52.56 vs. example 56.54) |
