package real

import (
	"testing"

	httpapi "akim5/api/internal/http"
)

func TestRealAdapterPreservesExplanationFacts(t *testing.T) {
	t.Setenv("DATA_DIR", "")
	decisions := []httpapi.Decision{
		{InitiativeID: "M7", DistrictID: "nura"},
		{InitiativeID: "M8", DistrictID: "nura"},
		{InitiativeID: "M10", DistrictID: "nura"},
		{InitiativeID: "M12"},
		{InitiativeID: "M5", DistrictID: "saryarka"},
	}
	result, violations := NewEngine().Simulate(decisions)
	if len(violations) != 0 || result == nil || result.Score == nil {
		t.Fatalf("valid scenario failed: result=%+v violations=%+v", result, violations)
	}
	llmResult := toScoringResult(result)
	if !llmResult.Submittable || len(llmResult.Districts) != len(result.Districts) || len(llmResult.Contributions) != len(decisions) {
		t.Fatalf("explanation lost scenario facts: districts=%d contributions=%d submittable=%t", len(llmResult.Districts), len(llmResult.Contributions), llmResult.Submittable)
	}
	for _, district := range llmResult.Districts {
		if district.ID == "nura" {
			if len(district.Indicators) != 10 {
				t.Fatalf("Nura indicators = %d, want 10", len(district.Indicators))
			}
			return
		}
	}
	t.Fatal("Nura district missing from explanation result")
}
