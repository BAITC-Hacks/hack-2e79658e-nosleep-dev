package httpapi

import (
	"context"

	"akim5/api/internal/events"
	"akim5/api/internal/llm"
	"akim5/api/internal/optimizer"
	"akim5/api/internal/scoring"
)

// NewLiveServices connects the HTTP contract to the production model.
func NewLiveServices() (Engine, Explainer, Deck, error) {
	model, err := scoring.New()
	if err != nil {
		return nil, nil, nil, err
	}
	deck, err := events.New()
	if err != nil {
		return nil, nil, nil, err
	}
	return &liveEngine{model: optimizer.New(model)}, &liveExplainer{client: llm.NewFromEnv()}, &liveDeck{deck: deck}, nil
}

type liveEngine struct{ model scoring.Engine }

func (e *liveEngine) Catalog() Catalog {
	source := e.model.Catalog()
	out := Catalog{Budget: source.Budget, RequiredDecisions: source.RequiredDecisions}
	for _, district := range source.Districts {
		out.Districts = append(out.Districts, District(district))
	}
	for _, initiative := range source.Initiatives {
		out.Initiatives = append(out.Initiatives, Initiative(initiative))
	}
	return out
}

func (e *liveEngine) Simulate(decisions []Decision) (*Result, []Violation) {
	result, violations := e.model.Simulate(toScoringDecisions(decisions))
	out := fromScoringResult(result)
	converted := make([]Violation, 0, len(violations))
	for _, violation := range violations {
		converted = append(converted, Violation{Code: violation.Code, Message: violation.Message})
	}
	return out, converted
}

func (e *liveEngine) Optimum(reveal bool) OptimumResult {
	source := e.model.Optimum(reveal)
	out := OptimumResult{BestScore: source.BestScore}
	for _, decision := range source.Decisions {
		out.Decisions = append(out.Decisions, Decision(decision))
	}
	return out
}

type liveExplainer struct{ client llm.Explainer }

func (e *liveExplainer) Explain(ctx context.Context, decisions []Decision, result *Result) Explanation {
	answer := e.client.Explain(ctx, toScoringDecisions(decisions), toScoringResult(result))
	return Explanation{Summary: answer.Summary, Strengths: answer.Strengths, Risks: answer.Risks, Recommendations: answer.Recommendations, Source: answer.Source}
}

func (e *liveExplainer) Compare(ctx context.Context, scenarios []NamedResult) Comparison {
	converted := make([]llm.NamedResult, 0, len(scenarios))
	for _, scenario := range scenarios {
		converted = append(converted, llm.NamedResult{Label: scenario.Label, Result: toScoringResult(scenario.Result)})
	}
	answer := e.client.Compare(ctx, converted)
	return Comparison{Summary: answer.Summary, Source: answer.Source}
}

type liveDeck struct{ deck events.Deck }

func (d *liveDeck) Draw(seed *int64) Event {
	return Event(d.deck.Draw(seed))
}

func toScoringDecisions(decisions []Decision) []scoring.Decision {
	converted := make([]scoring.Decision, 0, len(decisions))
	for _, decision := range decisions {
		converted = append(converted, scoring.Decision(decision))
	}
	return converted
}

func fromScoringResult(source *scoring.Result) *Result {
	if source == nil {
		return nil
	}
	out := &Result{
		BudgetUsed: source.BudgetUsed, BudgetRemaining: source.BudgetRemaining,
		DAvgBefore: source.DAvgBefore, DAvgAfter: source.DAvgAfter,
		MinDistrictID: source.MinDistrictID, MinDBefore: source.MinDBefore, MinDAfter: source.MinDAfter,
		CriticalBefore: source.CriticalBefore, CriticalAfter: source.CriticalAfter,
		Score: source.Score, Synergies: make([]Synergy, 0, len(source.Synergies)),
		Contributions: make([]Contribution, 0, len(source.Contributions)),
	}
	for _, district := range source.Districts {
		converted := DistrictResult{ID: district.ID, Name: district.Name, ScoreBefore: district.ScoreBefore, ScoreAfter: district.ScoreAfter}
		for _, indicator := range district.Indicators {
			converted.Indicators = append(converted.Indicators, IndicatorResult(indicator))
		}
		out.Districts = append(out.Districts, converted)
	}
	for _, label := range source.Synergies {
		out.Synergies = append(out.Synergies, Synergy{Label: label})
	}
	for _, contribution := range source.Contributions {
		out.Contributions = append(out.Contributions, Contribution(contribution))
	}
	return out
}

func toScoringResult(source *Result) *scoring.Result {
	if source == nil {
		return nil
	}
	out := &scoring.Result{
		BudgetUsed: source.BudgetUsed, BudgetRemaining: source.BudgetRemaining,
		DAvgBefore: source.DAvgBefore, DAvgAfter: source.DAvgAfter,
		MinDistrictID: source.MinDistrictID, MinDBefore: source.MinDBefore, MinDAfter: source.MinDAfter,
		CriticalBefore: source.CriticalBefore, CriticalAfter: source.CriticalAfter, Score: source.Score,
	}
	for _, district := range source.Districts {
		converted := scoring.DistrictResult{ID: district.ID, Name: district.Name, ScoreBefore: district.ScoreBefore, ScoreAfter: district.ScoreAfter}
		for _, indicator := range district.Indicators {
			converted.Indicators = append(converted.Indicators, scoring.Indicator(indicator))
		}
		out.Districts = append(out.Districts, converted)
	}
	for _, synergy := range source.Synergies {
		out.Synergies = append(out.Synergies, synergy.Label)
	}
	for _, contribution := range source.Contributions {
		out.Contributions = append(out.Contributions, scoring.Contribution(contribution))
	}
	return out
}
