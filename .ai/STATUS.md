# STATUS — demo readiness

_Updated 2026-09-23, Asia/Almaty. Keep this short for teammates and their agents._

**On `main`:** real Go scoring, catalog, explanation with labeled cached fallback, optimum, event deck, and web simulator with SVG map and scenario comparison are merged. PR #10 fixed local API startup and explanation data. Full Compose web startup still returns HTTP 500 because its image lacks root `contracts/` and `data/`.

**On `be/demo-readiness`:** `make dev` installs missing frontend dependencies on first run; the web image includes shared JSON assets; README quick start and feature status match the implementation.

**Verified:** clean-worktree `docker compose up --build -d` starts healthy API and web (HTTP 200); three browser rehearsals with `DEMO_MODE=true` and no key reached Score 56.54, a labeled cached explanation, optimum gauge, and two-scenario comparison. Fresh-worktree `make dev` installed dependencies and served both apps. Web image rebuild after local `npm ci` kept Docker context near 451 kB.

**Next:** get one teammate approval on PR #13, merge it, then tag and rehearse the merged `main` once more. The city-event endpoint has no web flow; keep it out of the 60-second demo.
