package httpapi

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FixtureServices implement the service boundary with the checked-in examples.
// They are intended for local demo startup until the AI/data services land.
type FixtureServices struct {
	catalog        Catalog
	baseline       *Result
	simulations    map[string]simulationFixture
	explainByKey   map[string]Explanation
	compareKey     string
	comparison     Comparison
	optimumDefault OptimumResult
	optimumReveal  OptimumResult
	event          Event
}

type simulationFixture struct {
	result     *Result
	violations []Violation
}

func NewFixtureServices(examplesDir string) (*FixtureServices, error) {
	examplesDir, err := resolveExamplesDir(examplesDir)
	if err != nil {
		return nil, err
	}
	service := &FixtureServices{
		simulations:  make(map[string]simulationFixture),
		explainByKey: make(map[string]Explanation),
	}
	if err := readFixture(filepath.Join(examplesDir, "catalog.json"), &service.catalog); err != nil {
		return nil, err
	}
	service.baseline = service.makeBaseline()
	service.loadBaseScore(examplesDir)
	if err := service.loadSimulation(filepath.Join(examplesDir, "simulate-valid.json")); err != nil {
		return nil, err
	}
	if err := service.loadSimulation(filepath.Join(examplesDir, "simulate-partial.json")); err != nil {
		return nil, err
	}
	if err := service.loadSimulation(filepath.Join(examplesDir, "simulate-invalid-budget.json")); err != nil {
		return nil, err
	}
	if err := service.loadExplain(filepath.Join(examplesDir, "explain-cached.json")); err != nil {
		return nil, err
	}
	if err := service.loadCompare(filepath.Join(examplesDir, "compare.json")); err != nil {
		return nil, err
	}
	if err := service.loadOptimum(filepath.Join(examplesDir, "optimum.json")); err != nil {
		return nil, err
	}
	if err := service.loadEvent(filepath.Join(examplesDir, "event.json")); err != nil {
		return nil, err
	}
	return service, nil
}

func (s *FixtureServices) Catalog() Catalog { return s.catalog }

func (s *FixtureServices) Simulate(decisions []Decision) (*Result, []Violation) {
	key := decisionsKey(decisions)
	if fixture, ok := s.simulations[key]; ok {
		return cloneResult(fixture.result), append([]Violation(nil), fixture.violations...)
	}

	result := cloneResult(s.baseline)
	result.BudgetUsed = 0
	for _, decision := range decisions {
		if initiative, ok := s.initiative(decision.InitiativeID); ok {
			result.BudgetUsed += initiative.Cost
		}
	}
	result.BudgetRemaining = s.catalog.Budget - result.BudgetUsed
	violations := s.validate(decisions)
	if len(violations) > 0 {
		result.Score = nil
	} else if len(decisions) == s.catalog.RequiredDecisions {
		// Unknown combinations use the baseline score as a visible placeholder;
		// known demo scenarios above always return their exact fixture numbers.
		value := 52.5577
		result.Score = &value
	}
	return result, violations
}

func (s *FixtureServices) Optimum(reveal bool) OptimumResult {
	if reveal {
		return s.optimumReveal
	}
	return s.optimumDefault
}

func (s *FixtureServices) Explain(_ context.Context, decisions []Decision, _ *Result) Explanation {
	if explanation, ok := s.explainByKey[decisionsKey(decisions)]; ok {
		return explanation
	}
	return Explanation{
		Summary:         "Сценарий рассчитан; подробное объяснение появится после подключения AI-сервиса.",
		Strengths:       []string{},
		Risks:           []string{},
		Recommendations: []string{},
		Source:          "cached",
	}
}

func (s *FixtureServices) Compare(_ context.Context, scenarios []NamedResult) Comparison {
	if namedResultsKey(scenarios) == s.compareKey {
		return s.comparison
	}
	return Comparison{Summary: "Сценарии рассчитаны; подробное сравнение появится после подключения AI-сервиса.", Source: "cached"}
}

func (s *FixtureServices) Draw(_ *int64) Event { return s.event }

func (s *FixtureServices) loadSimulation(path string) error {
	var fixture struct {
		Request struct {
			Decisions []Decision `json:"decisions"`
		} `json:"request"`
		Response json.RawMessage `json:"response"`
	}
	if err := readFixture(path, &fixture); err != nil {
		return err
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(fixture.Response, &fields); err != nil {
		return fmt.Errorf("decode %s response: %w", filepath.Base(path), err)
	}
	var violations []Violation
	if raw := fields["violations"]; len(raw) > 0 {
		if err := json.Unmarshal(raw, &violations); err != nil {
			return fmt.Errorf("decode %s violations: %w", filepath.Base(path), err)
		}
	}
	// The budget-overflow example intentionally uses explanatory placeholder
	// text for districts. Use the same baseline shape while retaining its exact
	// budget and violation values.
	placeholderDistricts := len(fields["districts"]) > 0 && fields["districts"][0] == '"'
	if placeholderDistricts {
		delete(fields, "districts")
	}
	clean, err := json.Marshal(fields)
	if err != nil {
		return err
	}
	result := Result{}
	if placeholderDistricts {
		result = *cloneResult(s.baseline)
	}
	if err := json.Unmarshal(clean, &result); err != nil {
		return fmt.Errorf("decode %s result: %w", filepath.Base(path), err)
	}
	if result.Districts == nil {
		result.Districts = s.makeBaseline().Districts
	}
	s.simulations[decisionsKey(fixture.Request.Decisions)] = simulationFixture{result: &result, violations: violations}
	return nil
}

func (s *FixtureServices) loadExplain(path string) error {
	var fixture struct {
		Request struct {
			Decisions []Decision `json:"decisions"`
		} `json:"request"`
		Response struct {
			Explanation Explanation `json:"explanation"`
		} `json:"response"`
	}
	if err := readFixture(path, &fixture); err != nil {
		return err
	}
	fixture.Response.Explanation.Source = "cached"
	s.explainByKey[decisionsKey(fixture.Request.Decisions)] = fixture.Response.Explanation
	return nil
}

func (s *FixtureServices) loadCompare(path string) error {
	var fixture struct {
		Request struct {
			Scenarios []struct {
				Label     string     `json:"label"`
				Decisions []Decision `json:"decisions"`
			} `json:"scenarios"`
		} `json:"request"`
		Response struct {
			Comparison Comparison `json:"comparison"`
		} `json:"response"`
	}
	if err := readFixture(path, &fixture); err != nil {
		return err
	}
	s.compareKey = decisionsLabelsKey(fixture.Request.Scenarios)
	s.comparison = fixture.Response.Comparison
	return nil
}

func (s *FixtureServices) loadOptimum(path string) error {
	var fixture struct {
		Default  OptimumResult `json:"default"`
		Revealed OptimumResult `json:"revealed"`
	}
	if err := readFixture(path, &fixture); err != nil {
		return err
	}
	s.optimumDefault = fixture.Default
	s.optimumReveal = fixture.Revealed
	return nil
}

func (s *FixtureServices) loadEvent(path string) error {
	var fixture struct {
		Response Event `json:"response"`
	}
	if err := readFixture(path, &fixture); err != nil {
		return err
	}
	s.event = fixture.Response
	return nil
}

func (s *FixtureServices) makeBaseline() *Result {
	result := &Result{
		BudgetRemaining: s.catalog.Budget,
		Districts:       make([]DistrictResult, 0, len(s.catalog.Districts)),
		DAvgBefore:      56.8624,
		DAvgAfter:       56.8624,
		MinDistrictID:   "nura",
		MinDBefore:      49.18,
		MinDAfter:       49.18,
		CriticalBefore:  2,
		CriticalAfter:   2,
		Synergies:       []Synergy{},
		Contributions:   []Contribution{},
	}
	for _, district := range s.catalog.Districts {
		score := baselineDistrictScore(district.ID)
		indicators := make([]IndicatorResult, 0, len(district.Indicators))
		for _, id := range []string{"T1", "T2", "E1", "E2", "S1", "S2", "B1", "B2", "C1", "C2"} {
			if value, ok := district.Indicators[id]; ok {
				indicators = append(indicators, IndicatorResult{Indicator: id, Before: value, After: value, Delta: 0})
			}
		}
		result.Districts = append(result.Districts, DistrictResult{ID: district.ID, Name: district.Name, ScoreBefore: score, ScoreAfter: score, Indicators: indicators})
	}
	return result
}

func baselineDistrictScore(id string) float64 {
	switch id {
	case "yesil":
		return 62.99
	case "almaty":
		return 57.06
	case "saryarka":
		return 54.65
	case "baikonur":
		return 56.63
	case "nura":
		return 49.18
	default:
		return 0
	}
}

func (s *FixtureServices) loadBaseScore(examplesDir string) {
	root := filepath.Dir(filepath.Dir(examplesDir))
	var fixture struct {
		Valid []struct {
			Name     string `json:"name"`
			Expected struct {
				Score float64 `json:"score"`
			} `json:"expected"`
		} `json:"valid"`
	}
	if err := readFixture(filepath.Join(root, "data", "fixtures", "golden.json"), &fixture); err != nil {
		return
	}
	for _, row := range fixture.Valid {
		if row.Name == "base_no_decisions" {
			value := row.Expected.Score
			s.baseline.Score = &value
			return
		}
	}
}

func (s *FixtureServices) initiative(id string) (Initiative, bool) {
	for _, initiative := range s.catalog.Initiatives {
		if initiative.ID == id {
			return initiative, true
		}
	}
	return Initiative{}, false
}

func (s *FixtureServices) districtKnown(id string) bool {
	for _, district := range s.catalog.Districts {
		if district.ID == id {
			return true
		}
	}
	return false
}

func (s *FixtureServices) validate(decisions []Decision) []Violation {
	violations := make([]Violation, 0)
	seen := make(map[string]bool)
	directions := make(map[string]int)
	cost := 0
	for _, decision := range decisions {
		initiative, ok := s.initiative(decision.InitiativeID)
		if !ok {
			violations = append(violations, Violation{Code: "UNKNOWN_INITIATIVE", Message: "неизвестное мероприятие: " + decision.InitiativeID})
			continue
		}
		if seen[initiative.ID] {
			violations = append(violations, Violation{Code: "DUPLICATE_INITIATIVE", Message: "мероприятие " + initiative.ID + " выбрано повторно"})
		}
		seen[initiative.ID] = true
		cost += initiative.Cost
		directions[initiative.Direction]++
		if initiative.Type == "city" && decision.DistrictID != "" {
			violations = append(violations, Violation{Code: "CITY_INITIATIVE_HAS_DISTRICT", Message: initiative.ID + " действует на весь город и не принимает districtId"})
		}
		if initiative.Type == "district" && decision.DistrictID == "" {
			violations = append(violations, Violation{Code: "DISTRICT_REQUIRED", Message: "для мероприятия " + initiative.ID + " нужно выбрать район"})
		}
		if decision.DistrictID != "" && !s.districtKnown(decision.DistrictID) {
			violations = append(violations, Violation{Code: "UNKNOWN_DISTRICT", Message: "неизвестный район: " + decision.DistrictID})
		}
	}
	if len(decisions) != s.catalog.RequiredDecisions {
		violations = append(violations, Violation{Code: "DECISION_COUNT", Message: fmt.Sprintf("нужно ровно %d решений, получено %d", s.catalog.RequiredDecisions, len(decisions))})
	}
	if cost > s.catalog.Budget {
		violations = append(violations, Violation{Code: "BUDGET_EXCEEDED", Message: fmt.Sprintf("бюджет превышен: %d из %d", cost, s.catalog.Budget)})
	}
	for _, direction := range []string{"transport", "ecology", "social", "safety", "services"} {
		count := directions[direction]
		if count > 2 {
			violations = append(violations, Violation{Code: "DIRECTION_LIMIT_EXCEEDED", Message: fmt.Sprintf("направление %s выбрано %d раз; максимум — 2", direction, count)})
		}
	}
	if seen["M1"] && seen["M3"] {
		violations = append(violations, Violation{Code: "INCOMPATIBLE", Message: "M1 и M3: либо BRT, либо ЛРТ, в любом районе", Details: map[string]any{"pair": []string{"M1", "M3"}}})
	}
	return violations
}

func decisionsKey(decisions []Decision) string {
	encoded, _ := json.Marshal(decisions)
	return string(encoded)
}

func namedResultsKey(scenarios []NamedResult) string {
	items := make([]struct {
		Label     string     `json:"label"`
		Decisions []Decision `json:"decisions"`
	}, 0, len(scenarios))
	for _, scenario := range scenarios {
		items = append(items, struct {
			Label     string     `json:"label"`
			Decisions []Decision `json:"decisions"`
		}{Label: scenario.Label, Decisions: scenario.Decisions})
	}
	encoded, _ := json.Marshal(items)
	return string(encoded)
}

func decisionsLabelsKey(scenarios []struct {
	Label     string     `json:"label"`
	Decisions []Decision `json:"decisions"`
}) string {
	encoded, _ := json.Marshal(scenarios)
	return string(encoded)
}

func cloneResult(result *Result) *Result {
	if result == nil {
		return nil
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return nil
	}
	var copy Result
	if json.Unmarshal(encoded, &copy) != nil {
		return nil
	}
	return &copy
}

func readFixture(path string, target any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read fixture %s: %w", path, err)
	}
	if err := json.Unmarshal(bytes, target); err != nil {
		return fmt.Errorf("decode fixture %s: %w", path, err)
	}
	return nil
}

func resolveExamplesDir(explicit string) (string, error) {
	candidates := []string{explicit, os.Getenv("CONTRACTS_DIR")}
	if explicit != "" {
		candidates = append(candidates, filepath.Join(explicit, "examples"))
	}
	if contractsDir := os.Getenv("CONTRACTS_DIR"); contractsDir != "" {
		candidates = append(candidates, filepath.Join(contractsDir, "examples"))
	}
	candidates = append(candidates,
		"contracts/examples",
		"../../contracts/examples",
		"../../../contracts/examples",
		"../../../../contracts/examples",
	)
	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		if _, err := os.Stat(filepath.Join(candidate, "catalog.json")); err == nil {
			absolute, err := filepath.Abs(candidate)
			if err == nil {
				return absolute, nil
			}
		}
	}
	return "", fmt.Errorf("could not find contracts/examples (set CONTRACTS_DIR)")
}
