DONE: scoring + LLM are merged to main; added exhaustive optimizer Engine wrapper for GET /optimum.
ON BRANCH / MERGED: ai/scoring / scoring+LLM merged, optimizer pending push/review.
TESTED HOW: GO111MODULE=off go test ./apps/api/internal/scoring (PASS); optimizer/LLM compile pending apps/api/go.mod.
BROKEN / RISKY: Go cannot go:embed root data/ from module apps/api; loader uses canonical root data/ and supports DATA_DIR.
NEEDS (from whom): Backend: create apps/api/go.mod; wire optimizer.New(scoring.MustNew()) and llm.NewFromEnv(); decide on mirrored embedded data assets.
NEXT: push optimizer; then add optional event deck or fix integration issues once module lands.
