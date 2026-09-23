# Go package boundary inside `apps/api` (BE ↔ AI+data)

No separate service — one Go binary, split by package so both people can work without touching each other's files. BE never edits `internal/scoring` or `internal/llm`; AI+data never edits `internal/http` or `cmd/`.

```
apps/api/
  cmd/server/main.go          BE — wiring, env loading, graceful shutdown
  internal/http/              BE — Gin routes, request/response DTOs, CORS, error middleware
  internal/scoring/           AI+data — the engine in contracts/scoring.md
  internal/llm/               AI+data — provider client, prompts/, DEMO_MODE fallback
  internal/optimizer/         AI+data — brute-force for GET /optimum
  internal/events/            AI+data — data/events.json loader + /events/draw
  internal/data/              AI+data — go:embed loaders for data/*.json (DATA_DIR override for tests)
```

`internal/http` depends on the other packages only through these interfaces (BE codes against them starting hour 1, using an in-memory fake, before AI+data's real implementation lands):

```go
package scoring

type Engine interface {
    Catalog() Catalog
    Simulate(decisions []Decision) (*Result, []Violation) // Result always populated (before/after), even when invalid
    Optimum(reveal bool) OptimumResult
}

package llm

type Explainer interface {
    Explain(ctx context.Context, decisions []Decision, result *scoring.Result) Explanation // Explanation.Source: "live"|"cached"
    Compare(ctx context.Context, scenarios []NamedResult) Comparison
}

package events

type Deck interface {
    Draw(seed *int64) Event
}
```

Exact field names for `Result`, `Violation`, `Explanation`, `OptimumResult`, `Event` follow `contracts/api.md`'s JSON shapes one-to-one (Go struct tags = the JSON keys shown there). If AI+data needs to change one of these shapes, they update `contracts/api.md` and this file in the same PR and tag BE for review — the old shape keeps working until that PR merges, per the hard rule that `main` stays runnable.
