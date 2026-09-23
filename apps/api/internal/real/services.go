// Package real adapts AI/data domain types to the HTTP DTO boundary.
package real

import (
	"context"

	"akim5/api/internal/events"
	httpapi "akim5/api/internal/http"
	"akim5/api/internal/llm"
	"akim5/api/internal/optimizer"
	"akim5/api/internal/scoring"
)

type Engine struct{ inner scoring.Engine }

func NewEngine() *Engine { return &Engine{inner: optimizer.New(scoring.MustNew())} }
func (e *Engine) Catalog() httpapi.Catalog {
	c := e.inner.Catalog()
	return httpapi.Catalog{Budget: c.Budget, RequiredDecisions: c.RequiredDecisions, Districts: districts(c.Districts), Initiatives: initiatives(c.Initiatives)}
}
func (e *Engine) Simulate(ds []httpapi.Decision) (*httpapi.Result, []httpapi.Violation) {
	r, v := e.inner.Simulate(toScoringDecisions(ds))
	return result(r), violations(v)
}
func (e *Engine) Optimum(reveal bool) httpapi.OptimumResult {
	r := e.inner.Optimum(reveal)
	return httpapi.OptimumResult{BestScore: r.BestScore, Decisions: decisions(r.Decisions)}
}

type Explainer struct{ inner llm.Explainer }

func NewExplainer() *Explainer { return &Explainer{inner: llm.NewFromEnv()} }
func (e *Explainer) Explain(ctx context.Context, ds []httpapi.Decision, r *httpapi.Result) httpapi.Explanation {
	a := e.inner.Explain(ctx, toScoringDecisions(ds), toScoringResult(r))
	return explanation(a)
}
func (e *Explainer) Compare(ctx context.Context, scenarios []httpapi.NamedResult) httpapi.Comparison {
	input := make([]llm.NamedResult, 0, len(scenarios))
	for _, s := range scenarios {
		input = append(input, llm.NamedResult{Label: s.Label, Result: toScoringResult(s.Result)})
	}
	a := e.inner.Compare(ctx, input)
	return httpapi.Comparison{Summary: a.Summary, Source: a.Source}
}

type Deck struct{ inner events.Deck }

func NewDeck() *Deck { return &Deck{inner: events.MustNew()} }
func (d *Deck) Draw(seed *int64) httpapi.Event {
	e := d.inner.Draw(seed)
	return httpapi.Event{ID: e.ID, Title: e.Title, Description: e.Description, DistrictID: e.DistrictID, Shocks: e.Shocks, BudgetDelta: e.BudgetDelta}
}

func toScoringDecisions(in []httpapi.Decision) []scoring.Decision {
	out := make([]scoring.Decision, len(in))
	for i, d := range in {
		out[i] = scoring.Decision{InitiativeID: d.InitiativeID, DistrictID: d.DistrictID}
	}
	return out
}
func decisions(in []scoring.Decision) []httpapi.Decision {
	out := make([]httpapi.Decision, len(in))
	for i, d := range in {
		out[i] = httpapi.Decision{InitiativeID: d.InitiativeID, DistrictID: d.DistrictID}
	}
	return out
}
func districts(in []scoring.District) []httpapi.District {
	out := make([]httpapi.District, len(in))
	for i, d := range in {
		out[i] = httpapi.District{ID: d.ID, Name: d.Name, Profile: d.Profile, Population: d.Population, Indicators: d.Indicators}
	}
	return out
}
func initiatives(in []scoring.Initiative) []httpapi.Initiative {
	out := make([]httpapi.Initiative, len(in))
	for i, d := range in {
		out[i] = httpapi.Initiative{ID: d.ID, Direction: d.Direction, Name: d.Name, Type: d.Type, Cost: d.Cost, Lag: d.Lag, Effects: d.Effects}
	}
	return out
}
func violations(in []scoring.Violation) []httpapi.Violation {
	out := make([]httpapi.Violation, len(in))
	for i, v := range in {
		out[i] = httpapi.Violation{Code: v.Code, Message: v.Message}
	}
	return out
}
func result(in *scoring.Result) *httpapi.Result {
	if in == nil {
		return nil
	}
	out := &httpapi.Result{BudgetUsed: in.BudgetUsed, BudgetRemaining: in.BudgetRemaining, DAvgBefore: in.DAvgBefore, DAvgAfter: in.DAvgAfter, MinDistrictID: in.MinDistrictID, MinDBefore: in.MinDBefore, MinDAfter: in.MinDAfter, CriticalBefore: in.CriticalBefore, CriticalAfter: in.CriticalAfter, Score: in.Score}
	out.Synergies = make([]httpapi.Synergy, len(in.Synergies))
	for i, s := range in.Synergies {
		out.Synergies[i] = httpapi.Synergy{Label: s.Label, DistrictID: s.DistrictID}
	}
	out.Contributions = make([]httpapi.Contribution, len(in.Contributions))
	for i, c := range in.Contributions {
		out.Contributions[i] = httpapi.Contribution{InitiativeID: c.InitiativeID, Name: c.Name, Direction: c.Direction, DistrictID: c.DistrictID, Cost: c.Cost, Lag: c.Lag, RealizedFraction: c.RealizedFraction, EffectsApplied: c.EffectsApplied}
	}
	out.Districts = make([]httpapi.DistrictResult, len(in.Districts))
	for i, d := range in.Districts {
		out.Districts[i] = httpapi.DistrictResult{ID: d.ID, Name: d.Name, ScoreBefore: d.ScoreBefore, ScoreAfter: d.ScoreAfter}
		for _, x := range d.Indicators {
			out.Districts[i].Indicators = append(out.Districts[i].Indicators, httpapi.IndicatorResult{Indicator: x.Indicator, Before: x.Before, After: x.After, Delta: x.Delta})
		}
	}
	return out
}
func toScoringResult(in *httpapi.Result) *scoring.Result {
	if in == nil {
		return nil
	}
	out := &scoring.Result{Submittable: in.Score != nil, BudgetUsed: in.BudgetUsed, BudgetRemaining: in.BudgetRemaining, DAvgBefore: in.DAvgBefore, DAvgAfter: in.DAvgAfter, MinDistrictID: in.MinDistrictID, MinDBefore: in.MinDBefore, MinDAfter: in.MinDAfter, CriticalBefore: in.CriticalBefore, CriticalAfter: in.CriticalAfter, Score: in.Score}
	out.Districts = make([]scoring.DistrictResult, len(in.Districts))
	for i, district := range in.Districts {
		out.Districts[i] = scoring.DistrictResult{ID: district.ID, Name: district.Name, ScoreBefore: district.ScoreBefore, ScoreAfter: district.ScoreAfter}
		out.Districts[i].Indicators = make([]scoring.Indicator, len(district.Indicators))
		for j, indicator := range district.Indicators {
			out.Districts[i].Indicators[j] = scoring.Indicator{Indicator: indicator.Indicator, Before: indicator.Before, After: indicator.After, Delta: indicator.Delta}
		}
	}
	for _, s := range in.Synergies {
		out.Synergies = append(out.Synergies, scoring.Synergy{Label: s.Label, DistrictID: s.DistrictID})
	}
	out.Contributions = make([]scoring.Contribution, len(in.Contributions))
	for i, contribution := range in.Contributions {
		out.Contributions[i] = scoring.Contribution{InitiativeID: contribution.InitiativeID, Name: contribution.Name, Direction: contribution.Direction, DistrictID: contribution.DistrictID, Cost: contribution.Cost, Lag: contribution.Lag, RealizedFraction: contribution.RealizedFraction, EffectsApplied: contribution.EffectsApplied}
	}
	return out
}
func explanation(in llm.Explanation) httpapi.Explanation {
	return httpapi.Explanation{Summary: in.Summary, Strengths: in.Strengths, Risks: in.Risks, Recommendations: in.Recommendations, Source: in.Source}
}
