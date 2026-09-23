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
type recommendationResponse struct {
	Recommendations []string `json:"recommendations"`
	allowed         []string
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
	base := strings.TrimSpace(os.Getenv("LLM_BASE_URL"))
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	model := strings.TrimSpace(os.Getenv("LLM_MODEL"))
	if model == "" {
		model = "gpt-4.1-mini"
	}
	return &Client{baseURL: strings.TrimRight(base, "/"), apiKey: strings.TrimSpace(os.Getenv("OPENAI_API_KEY")), model: model, demo: strings.EqualFold(os.Getenv("DEMO_MODE"), "true"), http: &http.Client{Timeout: timeout}}
}

func (c *Client) Explain(ctx context.Context, _ []scoring.Decision, result *scoring.Result) Explanation {
	fallback := cachedExplanation(result)
	if c.demo || c.apiKey == "" || c.model == "" || result == nil {
		return fallback
	}
	out := recommendationResponse{allowed: recommendationOptions(result)}
	payload := map[string]any{"verifiedFacts": explanationFacts(result), "allowedRecommendations": out.allowed}
	if c.request(ctx, recommendationSystemPrompt, payload, &out, func() bool { return validRecommendations(out.Recommendations, out.allowed) }) == nil {
		explanation := groundedExplanation(result)
		explanation.Recommendations = out.Recommendations
		explanation.Source = "live"
		return explanation
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
	if c.request(ctx, comparisonSystemPrompt, map[string]any{"verifiedFacts": comparisonFacts(scenarios)}, &out, func() bool { return validScenarioClaim(out.Summary) }) == nil {
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
	properties := map[string]any{}
	required := []string{}
	switch out.(type) {
	case *recommendationResponse:
		selection := out.(*recommendationResponse)
		properties["recommendations"] = map[string]any{"type": "array", "items": map[string]any{"type": "string", "enum": selection.allowed}, "minItems": 1, "maxItems": 2}
		required = append(required, "recommendations")
	default:
		properties["summary"] = map[string]string{"type": "string"}
		required = append(required, "summary")
	}
	return map[string]any{"type": "json_schema", "json_schema": map[string]any{"name": "city_explanation", "strict": true, "schema": map[string]any{"type": "object", "properties": properties, "required": required, "additionalProperties": false}}}
}
func validRecommendations(items, allowed []string) bool {
	if len(items) == 0 || len(items) > 2 {
		return false
	}
	options := make(map[string]bool, len(allowed))
	for _, item := range allowed {
		options[item] = true
	}
	seen := map[string]bool{}
	for _, item := range items {
		if !options[item] || seen[item] {
			return false
		}
		seen[item] = true
	}
	return true
}
func validScenarioClaim(s string) bool {
	if !validPlainText(s) {
		return false
	}
	lower := strings.ToLower(s)
	for _, phrase := range []string{"проблемы решены", "проблемы полностью решены", "проблем больше нет", "критические проблемы"} {
		if strings.Contains(lower, phrase) {
			return false
		}
	}
	for _, word := range strings.FieldsFunc(lower, func(r rune) bool { return !unicode.IsLetter(r) }) {
		for _, prefix := range []string{"внедрен", "внедрён", "внедря", "реализован", "реализу", "запущен", "построен", "создан", "устранен", "устранён", "эффектив", "активн", "значительн"} {
			if strings.HasPrefix(word, prefix) {
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
		case "один", "одна", "одно", "две", "два", "двух", "три", "трёх", "трех", "четыре", "четырёх", "четырех", "пять", "пяти":
			return false
		}
	}
	return true
}
func explanationFacts(r *scoring.Result) []string {
	facts := make([]string, 0, len(r.Contributions)+6)
	if r.DAvgAfter > r.DAvgBefore {
		facts = append(facts, "В расчётном сценарии взвешенный городской показатель улучшился.")
	} else if r.DAvgAfter < r.DAvgBefore {
		facts = append(facts, "В расчётном сценарии взвешенный городской показатель ухудшился.")
	} else {
		facts = append(facts, "В расчётном сценарии взвешенный городской показатель не изменился.")
	}
	if r.CriticalAfter == 0 {
		facts = append(facts, "В расчётном сценарии критических показателей не осталось.")
	} else if r.CriticalAfter < r.CriticalBefore {
		facts = append(facts, "В расчётном сценарии критических показателей стало меньше, но они ещё остались; их районы не указаны.")
	} else if r.CriticalAfter > r.CriticalBefore {
		facts = append(facts, "В расчётном сценарии критических показателей стало больше; их районы не указаны.")
	} else {
		facts = append(facts, "В расчётном сценарии число критических показателей не изменилось; их районы не указаны.")
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
		facts = append(facts, "Почти весь доступный бюджет распределён между выбранными инициативами.")
	}
	if len(r.Synergies) > 0 {
		facts = append(facts, "Расчёт учитывает синергию выбранных инициатив.")
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
func groundedExplanation(r *scoring.Result) Explanation {
	city := "средний показатель качества жизни в городе не изменился."
	strengths := []string{}
	if r.DAvgAfter > r.DAvgBefore {
		city = "средний показатель качества жизни в городе вырос."
		strengths = append(strengths, "Взвешенный городской показатель улучшился по расчёту.")
	} else if r.DAvgAfter < r.DAvgBefore {
		city = "средний показатель качества жизни в городе снизился."
	}
	critical := "Число показателей ниже критического порога не изменилось."
	if r.CriticalAfter == 0 {
		if r.CriticalBefore == 0 {
			critical = "Показателей ниже критического порога нет."
		} else {
			critical = "Показателей ниже критического порога больше нет."
		}
	} else if r.CriticalAfter < r.CriticalBefore {
		critical = "Показателей ниже критического порога стало меньше, но они ещё есть."
	} else if r.CriticalAfter > r.CriticalBefore {
		critical = "Показателей ниже критического порога стало больше."
	}
	if r.CriticalAfter < r.CriticalBefore {
		strengths = append(strengths, "Критических показателей стало меньше по расчёту.")
	}
	if len(r.Synergies) > 0 {
		strengths = append(strengths, "Расчёт учитывает синергию выбранных инициатив.")
	}
	if len(strengths) == 0 {
		strengths = append(strengths, "Выбранные инициативы включены в расчёт сценария.")
	}
	risks := []string{}
	if r.CriticalAfter > 0 {
		risks = append(risks, "В расчётном сценарии критические показатели ещё остаются.")
	}
	weakest := r.MinDistrictID
	for _, district := range r.Districts {
		if district.ID == r.MinDistrictID && district.Name != "" {
			weakest = district.Name
			break
		}
	}
	if weakest != "" {
		risks = append(risks, "Самый слабый по итоговому показателю район — "+weakest+".")
	}
	if len(risks) < 2 && r.BudgetUsed > 0 && r.BudgetRemaining >= 0 && r.BudgetRemaining*10 <= r.BudgetUsed+r.BudgetRemaining {
		risks = append(risks, "Почти весь доступный бюджет распределён между выбранными инициативами.")
	}
	if len(risks) == 0 {
		risks = append(risks, "Следует проверить показатели после расчётного сценария.")
	}
	return Explanation{Summary: "По расчёту " + city + " " + critical, Strengths: strengths, Risks: risks}
}
func recommendationOptions(r *scoring.Result) []string {
	options := []string{"Сравните расчётный результат с альтернативным набором инициатив."}
	weakest := r.MinDistrictID
	for _, district := range r.Districts {
		if district.ID == r.MinDistrictID && district.Name != "" {
			weakest = district.Name
			break
		}
	}
	if weakest != "" {
		options = append(options, "Проверьте показатели района "+weakest+" в альтернативном сценарии.")
	}
	if r.CriticalAfter > 0 {
		options = append(options, "Проверьте, какие показатели остаются критическими после расчёта.")
	}
	if r.BudgetUsed > 0 && r.BudgetRemaining >= 0 && r.BudgetRemaining*10 <= r.BudgetUsed+r.BudgetRemaining {
		options = append(options, "Сравните этот сценарий с вариантом, оставляющим больше бюджетного резерва.")
	}
	if len(r.Synergies) > 0 {
		options = append(options, "Оцените вклад синергии, сравнив сценарий с набором без неё.")
	}
	return options
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
