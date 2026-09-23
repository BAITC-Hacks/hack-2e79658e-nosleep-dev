# Kickoff prompt — Backend agent

Paste this into a fresh agent session opened in this repo.

---

You are the **backend** owner on a 3-person, no-lead hackathon team building «Аким на 5 часов», an AI city-budget simulator, in 5 hours. Read, in order: `.ai/WORKFLOW.md` (how we work — ownership, git, checkpoints, hard rules), `.ai/PLAN.md` (the product, the MVP, the task table), `contracts/api.md` (the endpoints you're implementing), `contracts/go-interfaces.md` (the exact boundary with AI+data's code — code against these Go interfaces, don't wait for their implementation).

You own `apps/api/cmd/**`, `apps/api/internal/http/**`, `Makefile`, `docker-compose.yml`, and any Dockerfiles. `internal/scoring`, `internal/llm`, `internal/optimizer`, `internal/events`, `internal/data`, and everything under `data/` belong to AI+data — read-only for you. If you need something to change there, write it under `NEEDS:` in `.ai/handoffs/backend.md`.

**Stack:** Go 1.26, module `akim5/api`, Gin, `gin-contrib/cors`. No database, no auth, no queues — this task doesn't need them.

**First hour:**
1. `go mod init akim5/api`, add Gin + cors. `GET /health` → `200 {"status":"ok"}`. CORS allowing `WEB_ORIGIN` (env, default `http://localhost:3000`).
2. Error middleware producing exactly the shape in `contracts/examples/error.json` for every non-2xx response.
3. Wire every endpoint in `contracts/api.md` against an **in-memory fake** implementing `scoring.Engine` and `llm.Explainer` (return `contracts/examples/*.json` verbatim, keyed by request shape) so the frontend can go live against you immediately, without waiting on AI+data.
4. Push. Update `.ai/handoffs/backend.md`.

**Then:** swap the fake for AI+data's real `internal/scoring`/`internal/llm` as they land (their handoff will say when) — this should be a one-line change in `cmd/server/main.go` if the interfaces in `contracts/go-interfaces.md` are respected on both sides. Add `Makefile` (`make dev` runs both apps), `docker-compose.yml` (api on 8000, web on 3000, both reading `.env`), and a couple of tests that matter: `/simulate` happy path, bad input (`DECISION_COUNT`), one full round-trip against a golden fixture from `data/fixtures/golden.json`.

Every external-facing endpoint needs input validation and a useful 4xx message — trust `internal/scoring`'s validation for the domain rules, but still guard against malformed JSON, wrong types, etc. at the HTTP layer.

Update `.ai/handoffs/backend.md` after each meaningful chunk (≤10 lines, overwrite don't append). Push at least every 30 minutes. Reply in under 15 lines with your understanding, your first 3 tasks with acceptance criteria, and any blocking question — then start immediately.
