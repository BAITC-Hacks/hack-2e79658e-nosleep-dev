package llm

import (
	"context"
	"testing"

	"akim5/api/internal/scoring"
)

func TestDemoModeNeverUsesNetworkAndIsClearlyCached(t *testing.T) {
	c := NewFromEnv()
	c.demo = true
	c.apiKey = ""
	r := &scoring.Result{BudgetUsed: 95, BudgetRemaining: 5, DAvgBefore: 56.8624, DAvgAfter: 58.0776, MinDistrictID: "nura", MinDBefore: 49.18, MinDAfter: 52.9625, CriticalBefore: 2, CriticalAfter: 0}
	s := 56.5431
	r.Score = &s
	got := c.Explain(context.Background(), nil, r)
	if got.Source != "cached" || got.Summary == "" || len(got.Strengths) == 0 || len(got.Risks) == 0 || len(got.Recommendations) == 0 {
		t.Fatalf("invalid cached explanation: %#v", got)
	}
}
