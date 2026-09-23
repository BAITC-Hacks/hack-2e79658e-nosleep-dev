// Package optimizer finds the best valid five-decision scenario once at startup.
package optimizer

import "akim5/api/internal/scoring"

// Engine decorates the scoring engine and fulfils scoring.Engine, adding a
// precomputed answer for GET /optimum while delegating catalog and simulation.
type Engine struct {
	base scoring.Engine
	best scoring.OptimumResult
}

func New(base scoring.Engine) *Engine {
	e := &Engine{base: base}
	e.best = search(base)
	return e
}

func (e *Engine) Catalog() scoring.Catalog { return e.base.Catalog() }
func (e *Engine) Simulate(ds []scoring.Decision) (*scoring.Result, []scoring.Violation) {
	return e.base.Simulate(ds)
}
func (e *Engine) Optimum(reveal bool) scoring.OptimumResult {
	if !reveal {
		return scoring.OptimumResult{BestScore: e.best.BestScore}
	}
	return cloneResult(e.best)
}

func cloneResult(in scoring.OptimumResult) scoring.OptimumResult {
	out := in
	out.Decisions = append([]scoring.Decision(nil), in.Decisions...)
	return out
}

// search exhaustively visits every combination of five unique initiatives and
// every legal district placement. The scoring engine remains the sole source
// of validation and score math, preventing optimizer/engine drift.
func search(base scoring.Engine) scoring.OptimumResult {
	catalog := base.Catalog()
	var best scoring.OptimumResult
	var choose func(start int, choices []scoring.Initiative)
	choose = func(start int, choices []scoring.Initiative) {
		if len(choices) == catalog.RequiredDecisions {
			evaluatePlacements(base, choices, func(ds []scoring.Decision, score float64) {
				if score > best.BestScore {
					best = scoring.OptimumResult{BestScore: score, Decisions: append([]scoring.Decision(nil), ds...)}
				}
			})
			return
		}
		for i := start; i <= len(catalog.Initiatives)-(catalog.RequiredDecisions-len(choices)); i++ {
			candidate := catalog.Initiatives[i]
			if overBudget(choices, candidate, catalog.Budget) {
				continue
			}
			choose(i+1, append(choices, candidate))
		}
	}
	choose(0, nil)
	return best
}

func overBudget(existing []scoring.Initiative, next scoring.Initiative, budget int) bool {
	cost := next.Cost
	for _, item := range existing {
		cost += item.Cost
	}
	return cost > budget
}
func evaluatePlacements(base scoring.Engine, items []scoring.Initiative, found func([]scoring.Decision, float64)) {
	districts := base.Catalog().Districts
	decisions := make([]scoring.Decision, len(items))
	var visit func(int)
	visit = func(index int) {
		if index == len(items) {
			result, violations := base.Simulate(decisions)
			if len(violations) == 0 && result.Score != nil {
				found(decisions, *result.Score)
			}
			return
		}
		initiative := items[index]
		if initiative.Type == "city" {
			decisions[index] = scoring.Decision{InitiativeID: initiative.ID}
			visit(index + 1)
			return
		}
		for _, district := range districts {
			decisions[index] = scoring.Decision{InitiativeID: initiative.ID, DistrictID: district.ID}
			visit(index + 1)
		}
	}
	visit(0)
}
