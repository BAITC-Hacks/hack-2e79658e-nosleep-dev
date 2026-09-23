DONE: scoring engine pushed (b9c86a5); LLM client added with JSON Schema, 2 attempts, 25s default timeout, stdout logs, DEMO_MODE/no-key cached response.
ON BRANCH / MERGED: ai/scoring / not yet merged.
TESTED HOW: GO111MODULE=off go test ./apps/api/internal/scoring (PASS); LLM compile pending apps/api/go.mod.
BROKEN / RISKY: Go cannot go:embed root data/ from module apps/api; loader uses canonical root data/ and supports DATA_DIR.
NEEDS (from whom): Backend: create apps/api/go.mod, wire scoring.New()/llm.NewFromEnv(), decide on mirrored embedded data assets.
NEXT: push LLM, then optimizer/events only after backend wiring confirms types.
