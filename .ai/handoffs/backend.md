DONE:
ON BRANCH / MERGED:
TESTED HOW:
BROKEN / RISKY:
NEEDS (from whom):
NEXT: scaffold apps/api (go.mod akim5/api, Gin, gin-contrib/cors). Add GET /health, CORS for WEB_ORIGIN, the error middleware producing contracts/examples/error.json's shape. Build internal/http against the scoring.Engine and llm.Explainer interfaces in contracts/go-interfaces.md, using an in-memory fake engine that returns contracts/examples/simulate-valid.json until AI+data's real one lands — don't block on them. See .ai/PLAN.md's task table and .ai/prompts/backend.md if pasting into a fresh agent.
