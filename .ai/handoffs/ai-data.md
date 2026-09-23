DONE: scoring + LLM merged; optimizer/events merged; scoring now embeds district/initiative/rule JSON and supports DATA_DIR only as a test override.
ON BRANCH / MERGED: ai/scoring / changes pending push/review.
TESTED HOW: GO111MODULE=off go test ./apps/api/internal/scoring and ./apps/api/internal/events (PASS); table test reads all current golden fixtures.
BROKEN / RISKY: root data and embedded Go assets must be kept in sync when authoring data; production no longer depends on working directory.
NEEDS (from whom): Backend: add apps/api/go.mod; wire optimizer.New(scoring.MustNew()), llm.NewFromEnv(), events.MustNew().
NEXT: full module/API integration test when backend pushes its scaffold.
