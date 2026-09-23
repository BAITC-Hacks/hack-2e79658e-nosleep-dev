package optimizer

import (
	"testing"

	"akim5/api/internal/scoring"
)

func TestOptimumIsAtLeastGoldenScenario(t *testing.T) {
	t.Setenv("DATA_DIR", "../../../../data")
	base, err := scoring.New()
	if err != nil {
		t.Fatal(err)
	}
	engine := New(base)
	got := engine.Optimum(true)
	if got.BestScore < 56.485102 {
		t.Fatalf("best score %.4f below verified fixture", got.BestScore)
	}
	if len(got.Decisions) != 5 {
		t.Fatalf("got %d decisions", len(got.Decisions))
	}
	if hidden := engine.Optimum(false); hidden.BestScore != got.BestScore || len(hidden.Decisions) != 0 {
		t.Fatalf("hidden optimum leaked decisions: %#v", hidden)
	}
}
