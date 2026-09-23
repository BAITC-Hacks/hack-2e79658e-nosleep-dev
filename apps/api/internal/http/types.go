package httpapi

import "context"

// These JSON types mirror contracts/api.md. They live in the HTTP package so
// the API can run against fixture services before the AI/data packages land.
type Decision struct {
	InitiativeID string `json:"initiativeId"`
	DistrictID   string `json:"districtId,omitempty"`
}

type Catalog struct {
	Budget            int          `json:"budget"`
	RequiredDecisions int          `json:"requiredDecisions"`
	Districts         []District   `json:"districts"`
	Initiatives       []Initiative `json:"initiatives"`
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

type Violation struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Details map[string]any `json:"details,omitempty"`
}

type Result struct {
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
	Synergies       []Synergy        `json:"synergies"`
	Contributions   []Contribution   `json:"contributions"`
}

type DistrictResult struct {
	ID          string            `json:"id"`
	Name        string            `json:"name,omitempty"`
	ScoreBefore float64           `json:"scoreBefore"`
	ScoreAfter  float64           `json:"scoreAfter"`
	Indicators  []IndicatorResult `json:"indicators,omitempty"`
}

type IndicatorResult struct {
	Indicator string  `json:"indicator"`
	Before    float64 `json:"before"`
	After     float64 `json:"after"`
	Delta     float64 `json:"delta"`
}

type Synergy struct {
	Label      string `json:"label"`
	DistrictID string `json:"districtId,omitempty"`
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

type OptimumResult struct {
	BestScore float64    `json:"bestScore"`
	Decisions []Decision `json:"decisions,omitempty"`
}

// Source is included at the top level of the /explain response, not inside
// the explanation object.
type Explanation struct {
	Summary         string   `json:"summary"`
	Strengths       []string `json:"strengths"`
	Risks           []string `json:"risks"`
	Recommendations []string `json:"recommendations"`
	Source          string   `json:"-"`
}

type NamedResult struct {
	Label      string
	Decisions  []Decision
	Result     *Result
	Violations []Violation
}

type Comparison struct {
	Summary string `json:"summary"`
	Source  string `json:"source"`
}

type Event struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DistrictID  string             `json:"districtId"`
	Shocks      map[string]float64 `json:"shocks"`
	BudgetDelta int                `json:"budgetDelta"`
}

// Engine, Explainer, and Deck are local mirrors of contracts/go-interfaces.md.
// The implementation can be swapped in cmd/server without coupling handlers
// to packages owned by AI/data.
type Engine interface {
	Catalog() Catalog
	Simulate(decisions []Decision) (*Result, []Violation)
	Optimum(reveal bool) OptimumResult
}

type Explainer interface {
	Explain(ctx context.Context, decisions []Decision, result *Result) Explanation
	Compare(ctx context.Context, scenarios []NamedResult) Comparison
}

type Deck interface {
	Draw(seed *int64) Event
}
