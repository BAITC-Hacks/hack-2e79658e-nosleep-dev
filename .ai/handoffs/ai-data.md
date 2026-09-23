DONE: scoring + LLM merged; optimizer PR #2 pushed; event deck + deterministic seeded draw added.
ON BRANCH / MERGED: ai/scoring / scoring+LLM merged; optimizer/events pending review.
TESTED HOW: GO111MODULE=off go test ./apps/api/internal/scoring and ./apps/api/internal/events (PASS); optimizer/LLM compile pending apps/api/go.mod.
BROKEN / RISKY: Go cannot go:embed root data/ from module apps/api; loader uses canonical root data/ and supports DATA_DIR.
NEEDS (from whom): Backend: create apps/api/go.mod; wire optimizer.New(scoring.MustNew()), llm.NewFromEnv(), events.MustNew(); decide on mirrored embedded data assets.
NEXT: push events and resolve full integration as soon as the backend module lands.
