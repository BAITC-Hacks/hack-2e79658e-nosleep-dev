// Package llm provides an OpenAI-compatible explainer with an honest fallback.
package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"akim5/api/internal/scoring"
)

type Explanation struct {
	Summary         string   `json:"summary"`
	Strengths       []string `json:"strengths"`
	Risks           []string `json:"risks"`
	Recommendations []string `json:"recommendations"`
	Source          string   `json:"-"`
}
type NamedResult struct {
	Label  string          `json:"label"`
	Result *scoring.Result `json:"result"`
}
type Comparison struct {
	Summary string `json:"summary"`
	Source  string `json:"source"`
}
type Explainer interface {
	Explain(context.Context, []scoring.Decision, *scoring.Result) Explanation
	Compare(context.Context, []NamedResult) Comparison
}
type Client struct {
	baseURL, apiKey, model string
	demo                   bool
	http                   *http.Client
}

func NewFromEnv() *Client {
	timeout := 25 * time.Second
	if n, err := strconv.Atoi(os.Getenv("LLM_TIMEOUT_SECONDS")); err == nil && n > 0 {
		timeout = time.Duration(n) * time.Second
	}
	base := os.Getenv("LLM_BASE_URL")
	if base == "" {
		base = "https://integrate.api.nvidia.com/v1"
	}
	return &Client{baseURL: strings.TrimRight(base, "/"), apiKey: os.Getenv("LLM_API_KEY"), model: os.Getenv("LLM_MODEL"), demo: strings.EqualFold(os.Getenv("DEMO_MODE"), "true"), http: &http.Client{Timeout: timeout}}
}

func (c *Client) Explain(ctx context.Context, _ []scoring.Decision, result *scoring.Result) Explanation {
	fallback := cachedExplanation(result)
	if c.demo || c.apiKey == "" || c.model == "" || result == nil {
		return fallback
	}
	payload := map[string]any{"verifiedFacts": explanationFacts(result)}
	var out Explanation
	if c.request(ctx, explanationSystemPrompt, payload, &out, func() bool { return validExplanation(out) }) == nil {
		out.Source = "live"
		return out
	}
	return fallback
}
func (c *Client) Compare(ctx context.Context, scenarios []NamedResult) Comparison {
	fallback := Comparison{Summary: cachedComparison(scenarios), Source: "cached"}
	if c.demo || c.apiKey == "" || c.model == "" {
		return fallback
	}
	var out struct {
		Summary string `json:"summary"`
	}
	if c.request(ctx, comparisonSystemPrompt, map[string]any{"verifiedFacts": comparisonFacts(scenarios)}, &out, func() bool { return validPlainText(out.Summary) }) == nil {
		return Comparison{Summary: out.Summary, Source: "live"}
	}
	return fallback
}

// request retries once if transport, provider, JSON, or shape validation fails.
func (c *Client) request(ctx context.Context, system string, payload any, out any, valid func() bool) error {
	var last error
	for attempt := 1; attempt <= 2; attempt++ {
		start := time.Now()
		err := c.requestOnce(ctx, system, payload, out)
		if err == nil && !valid() {
			err = fmt.Errorf("llm returned output outside required shape or grounded text rules")
		}
		log.Printf("llm call attempt=%d latency=%s success=%t", attempt, time.Since(start).Round(time.Millisecond), err == nil)
		if err == nil {
			return nil
		}
		last = err
	}
	return last
}
func (c *Client) requestOnce(ctx context.Context, system string, payload, out any) error {
	user, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	body, err := json.Marshal(map[string]any{"model": c.model, "temperature": 0.2, "messages": []map[string]string{{"role": "system", "content": system}, {"role": "user", "content": string(user)}}, "response_format": jsonSchema(out)})
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("llm status %s", resp.Status)
	}
	var wire struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&wire); err != nil {
		return err
	}
	if len(wire.Choices) == 0 {
		return fmt.Errorf("llm returned no choices")
	}
	return json.Unmarshal([]byte(wire.Choices[0].Message.Content), out)
}
func jsonSchema(out any) map[string]any {
	properties := map[string]any{"summary": map[string]string{"type": "string"}}
	required := []string{"summary"}
	if _, ok := out.(*Explanation); ok {
		items := map[string]string{"type": "string"}
		properties["strengths"] = map[string]any{"type": "array", "items": items}
		properties["risks"] = map[string]any{"type": "array", "items": items}
		properties["recommendations"] = map[string]any{"type": "array", "items": items}
		required = append(required, "strengths", "risks", "recommendations")
	}
	return map[string]any{"type": "json_schema", "json_schema": map[string]any{"name": "city_explanation", "strict": true, "schema": map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}}}
}
func validExplanation(e Explanation) bool {
	if !validPlainText(e.Summary) || len(e.Strengths) == 0 || len(e.Risks) == 0 || len(e.Recommendations) == 0 {
		return false
	}
	for _, group := range [][]string{e.Strengths, e.Risks, e.Recommendations} {
		for _, item := range group {
			if !validPlainText(item) {
				return false
			}
		}
	}
	return true
}
func validPlainText(s string) bool {
	if strings.TrimSpace(s) == "" {
		return false
	}
	for _, r := range s {
		if unicode.IsDigit(r) || r == '%' {
			return false
		}
	}
	for _, word := range strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) }) {
		switch word {
		case "один", "одна", "одно", "две", "два", "двух", "три", "трёх", "трех", "четыре", "четырёх", "четырех", "пять", "пяти", "несколько", "нескольких":
			return false
		}
	}
	return true
}
func explanationFacts(r *scoring.Result) []string {
	facts := []string{}
	if r.DAvgAfter > r.DAvgBefore {
		facts = append(facts, "Взвешенный городской показатель улучшился после решений.")
	} else if r.DAvgAfter < r.DAvgBefore {
		facts = append(facts, "Взвешенный городской показатель ухудшился после решений.")
	} else {
		facts = append(facts, "Взвешенный городской показатель не изменился после решений.")
	}
	if r.CriticalAfter == 0 {
		facts = append(facts, "После решений критических показателей не осталось.")
	} else if r.CriticalAfter < r.CriticalBefore {
		facts = append(facts, "Критических показателей стало меньше, но они ещё остались.")
	} else if r.CriticalAfter > r.CriticalBefore {
		facts = append(facts, "Критических показателей стало больше.")
	} else {
		facts = append(facts, "Число критических показателей не изменилось.")
	}
	names := map[string]string{}
	for _, district := range r.Districts {
		names[district.ID] = district.Name
	}
	weakest := names[r.MinDistrictID]
	if weakest == "" {
		weakest = r.MinDistrictID
	}
	if weakest != "" {
		facts = append(facts, "Самый слабый по итоговому районному показателю район — "+weakest+".")
	}
	if r.BudgetUsed > 0 && r.BudgetRemaining >= 0 && r.BudgetRemaining*10 <= r.BudgetUsed+r.BudgetRemaining {
		facts = append(facts, "Почти весь доступный бюджет использован.")
	}
	if len(r.Synergies) > 0 {
		facts = append(facts, "Сработала синергия выбранных инициатив.")
	}
	for _, item := range r.Contributions {
		if item.DistrictID == "" {
			facts = append(facts, "Выбрана общегородская инициатива: "+item.Name+".")
		} else {
			facts = append(facts, "Выбрана инициатива «"+item.Name+"» для района "+names[item.DistrictID]+".")
		}
	}
	return facts
}
func comparisonFacts(scenarios []NamedResult) []string {
	facts := make([]string, 0, len(scenarios)*2+1)
	var best *NamedResult
	for i := range scenarios {
		s := &scenarios[i]
		if s.Result == nil || s.Result.Score == nil {
			continue
		}
		if best == nil || *s.Result.Score > *best.Result.Score {
			best = s
		}
		if s.Result.DAvgAfter > s.Result.DAvgBefore {
			facts = append(facts, "В сценарии «"+s.Label+"» городской показатель улучшился.")
		} else if s.Result.DAvgAfter < s.Result.DAvgBefore {
			facts = append(facts, "В сценарии «"+s.Label+"» городской показатель ухудшился.")
		}
		if s.Result.CriticalAfter == 0 {
			facts = append(facts, "В сценарии «"+s.Label+"» критических показателей не осталось.")
		} else if s.Result.CriticalAfter < s.Result.CriticalBefore {
			facts = append(facts, "В сценарии «"+s.Label+"» критических показателей стало меньше, но они остались.")
		}
	}
	if best != nil {
		facts = append(facts, "По итоговому Score лидирует сценарий «"+best.Label+"».")
	}
	return facts
}
func cachedExplanation(r *scoring.Result) Explanation {
	if r == nil {
		return Explanation{Summary: "Нет результата для объяснения.", Strengths: []string{"Сценарий не рассчитан"}, Risks: []string{"Нужен корректный набор решений"}, Recommendations: []string{"Проверьте выбор"}, Source: "cached"}
	}
	score := "—"
	delta := "—"
	if r.Score != nil {
		score = fmt.Sprintf("%.2f", *r.Score)
		delta = fmt.Sprintf("%+.2f", *r.Score-(0.7*r.DAvgBefore+0.3*r.MinDBefore-float64(r.CriticalBefore)))
	}
	return Explanation{Summary: fmt.Sprintf("Расчётный Score: %s (%s к базовому сценарию).", score, delta), Strengths: []string{fmt.Sprintf("Средний показатель города: %.2f → %.2f.", r.DAvgBefore, r.DAvgAfter), fmt.Sprintf("Критических показателей: %d → %d.", r.CriticalBefore, r.CriticalAfter)}, Risks: []string{fmt.Sprintf("Самый слабый район после сценария: %s (%.2f).", r.MinDistrictID, r.MinDAfter), fmt.Sprintf("Использовано %d из %d единиц бюджета.", r.BudgetUsed, r.BudgetUsed+r.BudgetRemaining)}, Recommendations: []string{"Сохраняйте фокус на показателях ниже порога 40.", "Сравните сценарий с альтернативным набором инициатив."}, Source: "cached"}
}
func cachedComparison(s []NamedResult) string {
	if len(s) == 0 {
		return "Нет корректных сценариев для сравнения."
	}
	best := s[0]
	for _, x := range s[1:] {
		if x.Result != nil && x.Result.Score != nil && (best.Result == nil || best.Result.Score == nil || *x.Result.Score > *best.Result.Score) {
			best = x
		}
	}
	if best.Result == nil || best.Result.Score == nil {
		return "Нет корректных сценариев для сравнения."
	}
	return fmt.Sprintf("Сценарий %q лидирует по расчётному Score: %.2f.", best.Label, *best.Result.Score)
}
