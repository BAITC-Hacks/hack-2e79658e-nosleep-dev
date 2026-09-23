# STATUS

_Updated at every checkpoint by whoever merges last. Keep it short — this is read by teammates and their agents, not judges._

## T+0:00 — docs bootstrap
**What works on `main`:** nothing runnable yet. `.ai/`, `contracts/`, `data/`, `scripts/verify.py`, root config (`.gitignore`, `.env.example`) are pushed. No `apps/web` or `apps/api` code exists yet — each owner scaffolds their own per `.ai/handoffs/`.

**Blockers:** none. All three people can start in parallel right now.

**Next 60 min (→ Checkpoint 1):**
- FE: scaffold `apps/web` (Next.js + shadcn/ui), API client with `USE_MOCKS` reading `contracts/examples/*.json`.
- BE: scaffold `apps/api` (Gin, `/health`, CORS, error middleware) with an in-memory fake `scoring.Engine`.
- AI+data: start `internal/scoring` against `contracts/scoring.md`, aiming to pass `data/fixtures/golden.json` by end of hour 1.
