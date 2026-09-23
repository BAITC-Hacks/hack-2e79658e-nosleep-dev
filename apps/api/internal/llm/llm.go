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

func (c *Client) Explain(ctx context.Context, decisions []scoring.Decision, result *scoring.Result) Explanation {
	fallback := cachedExplanation(result)
	if c.demo || c.apiKey == "" || c.model == "" {
		return fallback
	}
	payload := map[string]any{"decisions": decisions, "result": result}
	var out Explanation
	if c.request(ctx, explanationSystemPrompt, payload, &out) == nil && validExplanation(out) {
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
	if c.request(ctx, comparisonSystemPrompt, map[string]any{"scenarios": scenarios}, &out) == nil && strings.TrimSpace(out.Summary) != "" {
		return Comparison{Summary: out.Summary, Source: "live"}
	}
	return fallback
}

// request retries once if transport, provider, JSON, or shape validation fails.
func (c *Client) request(ctx context.Context, system string, payload any, out any) error {
	var last error
	for attempt := 1; attempt <= 2; attempt++ {
		start := time.Now()
		err := c.requestOnce(ctx, system, payload, out)
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
	return strings.TrimSpace(e.Summary) != "" && len(e.Strengths) > 0 && len(e.Risks) > 0 && len(e.Recommendations) > 0
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
