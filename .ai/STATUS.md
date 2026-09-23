# STATUS — demo readiness

_Updated 2026-09-23, Asia/Almaty. Keep this short for teammates and their agents._

**On `main`:** real Go scoring, catalog, explanation with labeled cached fallback, optimum, event deck, and web simulator with SVG map and scenario comparison are merged. PR #13 fixed clean-clone startup and Compose packaging; PR #14 grounded AI output in verified scoring facts. PRs #16 and #17 switched the default provider to OpenAI and verified its configuration.

**On `ai/demo-qa`:** `DEMO.md` adds the 60-second show flow and pre-demo checks. No frontend or backend source files changed.

**Verified:** fresh remote clone at `ebca9de` installed dependencies via `make dev`; API `/health` and web both returned HTTP 200. Three browser rehearsals (desktop and narrow viewport) reached Score 56.54, critical 2 → 0, a labeled cached explanation, 98.8% optimum gauge, and two-scenario comparison. Over-budget UI blocked submit; saved scenarios survived reload. `docker compose up --build -d` in the same clone built both images, started a healthy API and web (HTTP 200), and returned the expected cached `/explain` result. An unreachable LLM also fell back to `source: cached` for `/explain` and `/compare` without changing Score. `go test ./...`, `python3 scripts/verify.py`, `npm run lint`, and `npm run build` all passed in the clone. Separately, the active local OpenAI `.env` returned `source: live` from `/explain` (Score 56.54307) and `/compare`; `go test ./...` passed after PR #17.

**Next:** review and merge this docs-only PR, then tag the final merged `main`. The frontend `/compare` request currently times out after 8s while the API allows a 25s LLM call; test live comparison latency before putting it in the demo. The city-event endpoint has no web flow; keep it out of the 60-second demo. Unsaved picks reset on reload; saved comparisons persist in browser storage.
