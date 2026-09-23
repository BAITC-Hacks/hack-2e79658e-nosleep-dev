# WORKFLOW — 3 peers, no lead, 5 hours

Adapted from the team's shared hackathon template for a 3-person, no-lead split. Read this once; it doesn't change during the event.

## Mission
Win with **one working, polished product that does one thing memorably well**, demoed flawlessly in under 60 seconds. Rank every choice: **judge-visible value > demo reliability > verifiability > everything else.** Cut scope before slipping time.

## Hard rules (never traded for speed)
- Real, honest commit history, made during the event. No fake or backdated commits.
- Everything judges need to verify works without our personal accounts: `DEMO_MODE` needs zero keys; any real key we use has a free tier or an organizer-provided path.
- Never present cached/canned output as live. `source: "live" | "cached"` is shown in the UI wherever the LLM speaks.
- `.env` is gitignored from commit 1. `.env.example` lists every variable. Never commit secrets.
- `main` is always runnable. If it's broken, fixing it beats every other task.

## Ownership
| Area | Owner | Others may |
|---|---|---|
| `apps/web/**` | Frontend | read only |
| `apps/api/cmd/**`, `apps/api/internal/http/**`, `Makefile`, `docker-compose.yml`, Dockerfiles | Backend | read only |
| `apps/api/internal/{scoring,llm,optimizer,events,data}/**`, `data/**` | AI+data | read only |
| `contracts/**` | shared — PR + approval from the *other* side of the contract | — |
| `.ai/PLAN.md`, `.ai/STATUS.md`, `.ai/DECISIONS.md`, `.ai/JUDGING.md`, root `README.md` | shared — whoever updates it, ping the other two | — |
| `.ai/handoffs/<person>.md` | that person | read |

Never edit a file you don't own without a PR. If you need a change there, write it under `NEEDS:` in your handoff and keep going with a local workaround (a stub, a mock, a `TODO`).

## Git workflow
- Branches: `fe/<topic>`, `be/<topic>`, `ai/<topic>`. Short-lived (≤45 min).
- Before starting: `git fetch && git rebase origin/main`.
- Before every commit: run tests/linters for what you touched, `git status`, stage explicit paths (never `git add .` blindly).
- Merge to `main`: open a PR, **one approval from either other teammate**, then merge. If both are heads-down, a self-merge after 15 minutes with a clear message is fine — don't block on ceremony.
- Commit messages describe the change: `feat(api): implement scoring engine per contracts/scoring.md`. Never `update`, `wip`, `hour 2`.
- Push at least every 30 minutes.

## Shared state
```
.ai/
  PLAN.md        problem, killer feature, MVP, non-goals, architecture, task list
  STATUS.md      what works on main right now, blockers, next 60 min — update at every checkpoint
  DECISIONS.md   one line per decision: what, why, rejected alternative
  JUDGING.md     criterion -> evidence in product -> evidence in demo -> weakness
  handoffs/
    frontend.md  backend.md  ai-data.md   — each person overwrites their own, ≤10 lines
contracts/       api.md, go-interfaces.md, scoring.md, examples/*.json
data/            districts.json, initiatives.json, rules.json, fixtures/golden.json (events.json to be added by AI+data)
```

Handoff format (overwrite after each meaningful chunk):
```text
DONE:
ON BRANCH / MERGED:
TESTED HOW:
BROKEN / RISKY:
NEEDS (from whom):
NEXT:
```

## Working rules
- Build the vertical slice first: pick → live `/simulate` call → map/Score update on screen. Ugly and complete beats pretty and partial.
- Stuck more than 10 minutes: stop, write it under `NEEDS:`/`BROKEN:` in your handoff, pick the simplest workaround, move on.
- Don't add dependencies you don't need, don't refactor working code, don't build infra beyond what's in the contract.
- Every user-facing flow needs loading, empty, error, and bad-input states before it counts as done.
- Every external call (LLM, GeoJSON fetch) has a timeout, error handling, and a fallback path.

## Timeline and checkpoints
| Time | Goal | Must be true on `main` at the end |
|---|---|---|
| T+0:00–0:15 | This docs push is the T+0 artifact | `.ai/PLAN.md`, `contracts/`, `data/` on `main`; both apps scaffold and boot |
| T+0:15–1:00 | Skeleton slice | Frontend calls a real backend endpoint and renders the (stubbed) result. One command starts everything |
| **Checkpoint 1** | | Tag `checkpoint-1`. Running slice + plan |
| T+1:00–2:00 | Real core logic | `/simulate` and `/catalog` are real, `internal/scoring` passes all golden fixtures. Stubs gone from the main path |
| **Checkpoint 2** | | Tag `checkpoint-2` |
| T+2:00–3:00 | Useful + robust | Error/loading/empty states, `/explain` live + DEMO_MODE fallback, AI output validated |
| **Checkpoint 3** | | Tag `checkpoint-3`. **Decide here whether the crisis event survives — cut it first if behind** |
| T+3:00–4:00 | Killer feature + polish | Live map visible in the demo path; UI polished; optimum gauge + compare if time allows |
| **Checkpoint 4** | | Tag `checkpoint-4` |
| T+4:00–4:15 | Last small fixes | — |
| **T+4:15 FEATURE FREEZE** | Only bug fixes, README, demo prep after this | |
| T+4:15–4:45 | README from a clean clone, demo rehearsal ×3, screenshots | A fresh clone following the README runs the demo |
| T+4:45–5:00 | Final merge, final tag, nothing risky | Tag `final`. Stop touching code |

Whoever merges last before a checkpoint tags `main` and updates `STATUS.md` — no fixed rotation, just don't let a checkpoint pass untagged.

## Cut order if behind (from `.ai/DECISIONS.md`)
1. City crisis event (`/events/draw`) — cut first if not solid by Checkpoint 3.
2. Scenario compare (`/compare`).
3. Optimum gauge (`/optimum`).
4. Real geo map → fall back to the stylized SVG (never cut the live-Score feature itself, only its rendering).
The 6 must-haves from the brief are never cut.
