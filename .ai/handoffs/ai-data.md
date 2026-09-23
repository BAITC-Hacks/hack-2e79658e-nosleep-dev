DONE: deterministic scoring engine + JSON-driven golden table test (3 valid / 9 invalid).
ON BRANCH / MERGED: ai/scoring / not yet pushed.
TESTED HOW: GO111MODULE=off go test ./apps/api/internal/scoring (PASS).
BROKEN / RISKY: Go cannot go:embed root data/ from module apps/api; loader uses canonical root data/ and supports DATA_DIR.
NEEDS (from whom): Backend: create apps/api/go.mod and wire scoring.New() as Engine; decide whether to mirror data under module for embedded release assets.
NEXT: commit/push scoring; then implement internal/llm with DEMO_MODE cached fallback.
