# PLAN — «Аким на 5 часов» / Astana Quality of Life simulator

**Problem:** a city manager choosing where to spend a limited budget across transport, ecology, social infrastructure, safety, and services has no fast way to see the tradeoffs of one scenario against another.
**User:** a city analyst or hackathon judge playing the role of an акim for one session.

**Product, in 3 sentences:** The player gets a fixed virtual budget (100 units) and picks exactly 5 initiatives from a 14-item catalog across 5 real Astana-style districts. A deterministic Go engine (never the LLM) computes how those picks move 10 quality-of-life indicators per district into a single **Astana Quality of Life Score**, live, as each pick is made. On submit, an LLM explains the result — strengths, risks, and tradeoffs — in plain language, never inventing numbers.

**KILLER FEATURE:** a live geo map of the 5 districts that recolors on every pick, with critical (<40) cells visibly turning from red to green as the Score ticks up in real time — visible within the first 60 seconds of the demo.

## MVP (must ship for the demo)
1. Single shared budget (100) and dataset, identical for every player — `data/*.json`, `GET /catalog`.
2. Exactly 5 decisions across the 5 directions, with the full rule set (budget, ≤2/direction, incompatibilities, synergies) enforced server-side — `POST /simulate`.
3. Live map + Score that update on every pick, before the set is even valid.
4. Final Score computed on submit, with a clear AI explanation of strengths/risks/tradeoffs — `POST /explain`.
5. `DEMO_MODE` fallback so the demo survives a dead network or key.
6. Two extras that double as judged "optional" points: an **optimum gauge** (`GET /optimum`, "your X = Y% of best possible") and **A/B/C scenario compare** (`POST /compare`, client-side storage, no DB).

## Non-goals (explicitly cut, not forgotten)
- Auth, accounts, multi-team persistence, a database — nothing here needs one at this scope.
- A live team leaderboard (needs shared persistence — cut).
- Hosted deploy — the demo runs local via `docker compose up` / `make dev`.
- AI-generated pitch deck.
- AI-proposed "best swap" recommendations beyond what `/optimum` already covers.
- Jev (TypeSafe's System One model) — released 2026-09-15, early access, returns typed choices/scores with no prose, not text explanations. Worth revisiting post-hackathon; not this build. Noted under "what's next" in the README.
- A city crisis event is in scope (`POST /events/draw`) but is the **first thing cut** if the team falls behind — see the cut order in `.ai/DECISIONS.md`.

## Stack & architecture
- **Backend:** Go 1.26 + Gin, one binary (`apps/api`), split into packages by owner (`contracts/go-interfaces.md`). No database — `data/*.json` is embedded via `go:embed` and is the single source of truth for both the engine and the frontend catalog.
- **Frontend:** Next.js (App Router, TypeScript), shadcn/ui + Tailwind as the base, 2–3 ObsidianUI components for the Score reveal and hero. A real geo map (Leaflet/MapLibre, OSM-sourced GeoJSON) with a stylized-SVG fallback if the boundaries aren't clean.
- **AI:** one OpenAI-compatible client (`internal/llm`), NVIDIA NIM by default, switchable to OpenAI via env. Structured output only (tool calling / JSON schema). Called on submit and on compare — never on every pick, never from the frontend.
- **No lead.** Three peers (FE / BE / AI+data), PR + one-approval merges to `main`. Details in `.ai/WORKFLOW.md`.

## Task list

| Owner | Task | Acceptance criteria | Due |
|---|---|---|---|
| BE | Scaffold `apps/api`: Gin, `/health`, CORS, error middleware, in-memory fake `scoring.Engine` | `curl localhost:8000/health` → 200; `go run ./cmd/server` boots | Checkpoint 1 |
| FE | Scaffold `apps/web`: Next.js, shadcn/ui installed, API client with `USE_MOCKS` switch reading `contracts/examples/*.json` | Landing page renders, calls `/health` (mocked or real) | Checkpoint 1 |
| AI+data | Implement `internal/scoring` exactly per `contracts/scoring.md`, passing every case in `data/fixtures/golden.json` | `go test ./internal/scoring/...` green on all golden fixtures | Checkpoint 1–2 |
| BE | Wire `POST /simulate`, `GET /catalog` to the real engine | `contracts/examples/simulate-valid.json` reproduced byte-for-byte (modulo float rounding) | Checkpoint 2 |
| FE | Build the decision flow: catalog → 5 picks → live map + Score, using `/simulate` | Picking any initiative updates the map and Score within one request round-trip | Checkpoint 2 |
| AI+data | Implement `internal/llm` (NVIDIA NIM client, structured output, DEMO_MODE, 1 retry then cached fallback) | `POST /explain` with `DEMO_MODE=true` and no key returns a valid cached response | Checkpoint 2–3 |
| FE | Submit flow: loading/empty/error/bad-input states, AI explanation panel with the `cached` badge | All 4 states demoed manually | Checkpoint 3 |
| AI+data | `internal/optimizer` (`GET /optimum`) and event deck (`internal/events`, `data/events.json`) if time allows | `/optimum` returns a score ≥ every fixture's score | Checkpoint 3–4 |
| FE | Compare UI (2–3 saved scenarios, browser storage) | Compare view renders 2 real scenarios side by side | Checkpoint 4 |
| Whoever's free | README quick start verified from a clean clone; demo script; 3 rehearsals | Fresh clone → documented commands → working demo | T+4:15–4:45 |

See `.ai/handoffs/` for the first concrete task each person should open a branch for right now.
