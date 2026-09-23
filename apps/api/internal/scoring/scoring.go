// Package scoring implements the deterministic city-budget scoring contract.
package scoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Decision struct {
	InitiativeID string `json:"initiativeId"`
	DistrictID   string `json:"districtId,omitempty"`
}
type Violation struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Indicator struct {
	Indicator string  `json:"indicator"`
	Before    float64 `json:"before"`
	After     float64 `json:"after"`
	Delta     float64 `json:"delta"`
}
type DistrictResult struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	ScoreBefore float64     `json:"scoreBefore"`
	ScoreAfter  float64     `json:"scoreAfter"`
	Indicators  []Indicator `json:"indicators"`
}
type Contribution struct {
	InitiativeID     string             `json:"initiativeId"`
	Name             string             `json:"name"`
	Direction        string             `json:"direction"`
	DistrictID       string             `json:"districtId,omitempty"`
	Cost             int                `json:"cost"`
	Lag              int                `json:"lag"`
	RealizedFraction float64            `json:"realizedFraction"`
	EffectsApplied   map[string]float64 `json:"effectsApplied"`
}
type Result struct {
	Submittable     bool             `json:"submittable"`
	Violations      []Violation      `json:"violations,omitempty"`
	BudgetUsed      int              `json:"budgetUsed"`
	BudgetRemaining int              `json:"budgetRemaining"`
	Districts       []DistrictResult `json:"districts"`
	DAvgBefore      float64          `json:"dAvgBefore"`
	DAvgAfter       float64          `json:"dAvgAfter"`
	MinDistrictID   string           `json:"minDistrictId"`
	MinDBefore      float64          `json:"minDBefore"`
	MinDAfter       float64          `json:"minDAfter"`
	CriticalBefore  int              `json:"criticalBefore"`
	CriticalAfter   int              `json:"criticalAfter"`
	Score           *float64         `json:"score"`
	Synergies       []string         `json:"synergies"`
	Contributions   []Contribution   `json:"contributions"`
}
type Catalog struct {
	Budget            int          `json:"budget"`
	RequiredDecisions int          `json:"requiredDecisions"`
	Districts         []District   `json:"districts"`
	Initiatives       []Initiative `json:"initiatives"`
}
type OptimumResult struct {
	BestScore float64    `json:"bestScore"`
	Decisions []Decision `json:"decisions,omitempty"`
}
type Engine interface {
	Catalog() Catalog
	Simulate([]Decision) (*Result, []Violation)
	Optimum(bool) OptimumResult
}
type District struct {
	ID         string             `json:"id"`
	Name       string             `json:"name"`
	Profile    string             `json:"profile"`
	Population float64            `json:"population"`
	Indicators map[string]float64 `json:"indicators"`
}
type Initiative struct {
	ID        string             `json:"id"`
	Direction string             `json:"direction"`
	Name      string             `json:"name"`
	Type      string             `json:"type"`
	Cost      int                `json:"cost"`
	Lag       int                `json:"lag"`
	Effects   map[string]float64 `json:"effects"`
}
type rules struct {
	Budget, HorizonQuarters, RequiredDecisions, MaxPerDirection int
	CriticalThreshold, CriticalPenaltyPerPair                   float64
	ScoreFormula                                                struct{ DAvgWeight, MinDistrictWeight float64 }
	IndicatorWeights                                            map[string]float64
	Synergies                                                   []synergy
	Incompatibilities                                           []incompatibility
}
type synergy struct {
	Pair          []string
	Indicator     string
	Bonus         float64
	Anchor, Label string
}
type incompatibility struct {
	Pair             []string
	SameDistrictOnly bool
	Label            string
}

type engine struct {
	catalog     Catalog
	rules       rules
	initiatives map[string]Initiative
	districts   map[string]District
}

// New loads the repository's canonical JSON data. DATA_DIR may point to an
// alternate data directory in integration tests or deployments.
func New() (*engine, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}
	var districts []District
	var initiatives []Initiative
	var r rules
	if err := decode(filepath.Join(dir, "districts.json"), &districts); err != nil {
		return nil, err
	}
	if err := decode(filepath.Join(dir, "initiatives.json"), &initiatives); err != nil {
		return nil, err
	}
	if err := decode(filepath.Join(dir, "rules.json"), &r); err != nil {
		return nil, err
	}
	e := &engine{catalog: Catalog{Budget: r.Budget, RequiredDecisions: r.RequiredDecisions, Districts: districts, Initiatives: initiatives}, rules: r, initiatives: map[string]Initiative{}, districts: map[string]District{}}
	for _, v := range initiatives {
		e.initiatives[v.ID] = v
	}
	for _, v := range districts {
		e.districts[v.ID] = v
	}
	return e, nil
}
func MustNew() *engine {
	e, err := New()
	if err != nil {
		panic(err)
	}
	return e
}
func (e *engine) Catalog() Catalog           { return e.catalog }
func (e *engine) Optimum(bool) OptimumResult { return OptimumResult{} }
func decode(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read %s: %w", path, err)
	}
	return json.Unmarshal(b, out)
}
func dataDir() (string, error) {
	if d := os.Getenv("DATA_DIR"); d != "" {
		return d, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(wd, "data", "rules.json")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Dir(candidate), nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "", fmt.Errorf("data directory not found; set DATA_DIR")
}

func (e *engine) Simulate(decisions []Decision) (*Result, []Violation) {
	v := e.validate(decisions)
	r := e.calculate(decisions)
	r.Violations = v
	r.Submittable = len(v) == 0
	if r.Submittable {
		s := score(r.DAvgAfter, r.MinDAfter, r.CriticalAfter, e.rules)
		r.Score = &s
	}
	return r, v
}
func (e *engine) validate(ds []Decision) []Violation {
	if len(ds) != e.rules.RequiredDecisions {
		return []Violation{{"DECISION_COUNT", fmt.Sprintf("нужно ровно %d решений, получено %d", e.rules.RequiredDecisions, len(ds))}}
	}
	seen := map[string]bool{}
	for _, d := range ds {
		if _, ok := e.initiatives[d.InitiativeID]; !ok {
			return []Violation{{"UNKNOWN_INITIATIVE", fmt.Sprintf("неизвестная инициатива: %s", d.InitiativeID)}}
		}
		if seen[d.InitiativeID] {
			return []Violation{{"DUPLICATE_INITIATIVE", fmt.Sprintf("инициатива повторяется: %s", d.InitiativeID)}}
		}
		seen[d.InitiativeID] = true
	}
	for _, d := range ds {
		in := e.initiatives[d.InitiativeID]
		if in.Type == "district" {
			if d.DistrictID == "" {
				return []Violation{{"DISTRICT_REQUIRED", fmt.Sprintf("для %s нужен район", d.InitiativeID)}}
			}
			if _, ok := e.districts[d.DistrictID]; !ok {
				return []Violation{{"UNKNOWN_DISTRICT", fmt.Sprintf("неизвестный район: %s", d.DistrictID)}}
			}
		} else if d.DistrictID != "" {
			return []Violation{{"CITY_INITIATIVE_HAS_DISTRICT", fmt.Sprintf("городская инициатива %s не принимает район", d.InitiativeID)}}
		}
	}
	cost := 0
	directions := map[string]int{}
	for _, d := range ds {
		in := e.initiatives[d.InitiativeID]
		cost += in.Cost
		directions[in.Direction]++
	}
	if cost > e.rules.Budget {
		return []Violation{{"BUDGET_EXCEEDED", fmt.Sprintf("бюджет превышен: %d из %d", cost, e.rules.Budget)}}
	}
	for direction, n := range directions {
		if n > e.rules.MaxPerDirection {
			return []Violation{{"DIRECTION_LIMIT_EXCEEDED", fmt.Sprintf("слишком много инициатив направления %s", direction)}}
		}
	}
	byID := map[string]Decision{}
	for _, d := range ds {
		byID[d.InitiativeID] = d
	}
	for _, x := range e.rules.Incompatibilities {
		a, aok := byID[x.Pair[0]]
		b, bok := byID[x.Pair[1]]
		if aok && bok && (!x.SameDistrictOnly || a.DistrictID == b.DistrictID) {
			return []Violation{{"INCOMPATIBLE", x.Label}}
		}
	}
	return nil
}
func (e *engine) calculate(ds []Decision) *Result {
	deltas := map[string]map[string]float64{}
	for id := range e.districts {
		deltas[id] = map[string]float64{}
	}
	cost := 0
	contributions := []Contribution{}
	selected := map[string]Decision{}
	for _, d := range ds {
		in, ok := e.initiatives[d.InitiativeID]
		if !ok {
			continue
		}
		selected[in.ID] = d
		cost += in.Cost
		fraction := float64(e.rules.HorizonQuarters-in.Lag) / float64(e.rules.HorizonQuarters)
		applied := map[string]float64{}
		for k, v := range in.Effects {
			applied[k] = v * fraction
			if in.Type == "city" {
				for id := range e.districts {
					deltas[id][k] += applied[k]
				}
			} else if _, ok := e.districts[d.DistrictID]; ok {
				deltas[d.DistrictID][k] += applied[k]
			}
		}
		contributions = append(contributions, Contribution{in.ID, in.Name, in.Direction, d.DistrictID, in.Cost, in.Lag, fraction, applied})
	}
	synergies := []string{}
	for _, s := range e.rules.Synergies {
		_, a := selected[s.Pair[0]]
		_, b := selected[s.Pair[1]]
		if a && b {
			anchor := selected[s.Anchor]
			if _, ok := e.districts[anchor.DistrictID]; ok {
				deltas[anchor.DistrictID][s.Indicator] += s.Bonus
				synergies = append(synergies, fmt.Sprintf("%s+%s", s.Pair[0], s.Pair[1]))
			}
		}
	}
	result := &Result{BudgetUsed: cost, BudgetRemaining: e.rules.Budget - cost, Synergies: synergies, Contributions: contributions}
	var avgBefore, avgAfter float64
	minBefore, minAfter := 101.0, 101.0
	for _, d := range e.catalog.Districts {
		beforeD, afterD := 0.0, 0.0
		indicators := make([]Indicator, 0, len(e.rules.IndicatorWeights))
		for k, w := range e.rules.IndicatorWeights {
			before := d.Indicators[k]
			after := clip(before + deltas[d.ID][k])
			indicators = append(indicators, Indicator{k, before, after, after - before})
			beforeD += w * before
			afterD += w * after
			if before < e.rules.CriticalThreshold {
				result.CriticalBefore++
			}
			if after < e.rules.CriticalThreshold {
				result.CriticalAfter++
			}
		}
		avgBefore += d.Population * beforeD
		avgAfter += d.Population * afterD
		if beforeD < minBefore {
			minBefore = beforeD
		}
		if afterD < minAfter {
			minAfter = afterD
			result.MinDistrictID = d.ID
		}
		result.Districts = append(result.Districts, DistrictResult{d.ID, d.Name, beforeD, afterD, indicators})
	}
	result.DAvgBefore, result.DAvgAfter, result.MinDBefore, result.MinDAfter = avgBefore, avgAfter, minBefore, minAfter
	return result
}
func clip(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
func score(avg, min float64, crit int, r rules) float64 {
	return r.ScoreFormula.DAvgWeight*avg + r.ScoreFormula.MinDistrictWeight*min - r.CriticalPenaltyPerPair*float64(crit)
}
