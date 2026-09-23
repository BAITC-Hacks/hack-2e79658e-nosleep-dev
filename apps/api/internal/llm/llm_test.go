package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"akim5/api/internal/scoring"
)

func TestNewFromEnvUsesOpenAIDefaults(t *testing.T) {
	t.Setenv("LLM_BASE_URL", "")
	t.Setenv("LLM_MODEL", "")
	t.Setenv("OPENAI_API_KEY", "test-key")
	t.Setenv("DEMO_MODE", "false")

	c := NewFromEnv()
	if c.baseURL != "https://api.openai.com/v1" || c.model != "gpt-4.1-mini" || c.apiKey != "test-key" || c.demo {
		t.Fatalf("unexpected OpenAI defaults: baseURL=%q model=%q keySet=%t demo=%t", c.baseURL, c.model, c.apiKey != "", c.demo)
	}
}

func TestDemoModeNeverUsesNetworkAndIsClearlyCached(t *testing.T) {
	c := NewFromEnv()
	c.demo = true
	c.apiKey = ""
	r := &scoring.Result{BudgetUsed: 95, BudgetRemaining: 5, DAvgBefore: 56.8624, DAvgAfter: 58.0776, MinDistrictID: "nura", MinDBefore: 49.18, MinDAfter: 52.9625, CriticalBefore: 2, CriticalAfter: 0}
	s := 56.5431
	r.Score = &s
	got := c.Explain(context.Background(), nil, r)
	if got.Source != "cached" || got.Summary == "" || len(got.Strengths) == 0 || len(got.Risks) == 0 || len(got.Recommendations) == 0 {
		t.Fatalf("invalid cached explanation: %#v", got)
	}
}

func TestExplainRetriesUngroundedText(t *testing.T) {
	requests := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Messages) != 2 || !strings.Contains(body.Messages[1].Content, "verifiedFacts") || strings.Contains(body.Messages[1].Content, "scoreAfter") {
			t.Fatalf("provider received ungrounded input: %#v", body.Messages)
		}
		answer := `{"recommendations":["Показатель города вырос на 3.78%."]}`
		if requests == 2 {
			answer = `{"recommendations":["Проверьте показатели района Нура в альтернативном сценарии."]}`
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": answer}}}})
	}))
	defer provider.Close()
	c := &Client{baseURL: provider.URL, apiKey: "test", model: "test", http: provider.Client()}
	got := c.Explain(context.Background(), nil, testResult())
	if requests != 2 || got.Source != "live" || strings.Contains(got.Summary, "%") || !strings.Contains(got.Summary, "критических показателей не осталось") {
		t.Fatalf("want grounded live result after retry; requests=%d result=%#v", requests, got)
	}
}

func TestExplainFallsBackAfterUngroundedRetries(t *testing.T) {
	requests := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"recommendations\":[\"Рост на 3.78%\"]}"}}]}`))
	}))
	defer provider.Close()
	c := &Client{baseURL: provider.URL, apiKey: "test", model: "test", http: provider.Client()}
	got := c.Explain(context.Background(), nil, testResult())
	if requests != 2 || got.Source != "cached" {
		t.Fatalf("want honest fallback after retry; requests=%d result=%#v", requests, got)
	}
}

func TestExplanationFactsCriticalPairsNotDistricts(t *testing.T) {
	facts := strings.Join(explanationFacts(testResult()), " ")
	if !strings.Contains(facts, "критических показателей не осталось") || strings.Contains(facts, "районах остались") {
		t.Fatalf("incorrect critical-pair description: %s", facts)
	}
}

func TestGroundedExplanationDoesNotAssignCriticalIndicatorsToWeakestDistrict(t *testing.T) {
	r := testResult()
	r.CriticalAfter = 1
	got := groundedExplanation(r)
	if !strings.Contains(got.Summary, "критических показателей стало меньше") || strings.Contains(got.Summary, "Нура") {
		t.Fatalf("summary must describe the calculation without inventing a district: %q", got.Summary)
	}
	if len(got.Risks) != 2 || !strings.Contains(got.Risks[0], "критические показатели") || !strings.Contains(got.Risks[1], "Нура") {
		t.Fatalf("risks must keep critical indicators and weakest district separate: %#v", got.Risks)
	}
}

func TestRecommendationOptionsKeepCriticalIndicatorsSeparateFromDistrict(t *testing.T) {
	r := testResult()
	r.CriticalAfter = 1
	options := recommendationOptions(r)
	joined := strings.Join(options, " ")
	if !strings.Contains(joined, "критическими") || !strings.Contains(joined, "Нура") {
		t.Fatalf("expected scenario options for remaining critical indicators and weakest district: %#v", options)
	}
	for _, option := range options {
		if strings.Contains(option, "критическими") && strings.Contains(option, "Нура") {
			t.Fatalf("option must not assign critical indicators to a district: %q", option)
		}
	}
}

func TestValidPlainTextAllowsVagueQuantity(t *testing.T) {
	if !validPlainText("Для района выбраны несколько инициатив.") {
		t.Fatal("a vague count of selected initiatives should not force a cached fallback")
	}
	if validPlainText("Выбраны пять инициатив.") {
		t.Fatal("exact counts remain excluded from AI prose")
	}
}

func TestRecommendationSchemaRequiresOneOrTwoItems(t *testing.T) {
	jsonSchemaBody := jsonSchema(&recommendationResponse{allowed: []string{"A", "B"}})["json_schema"].(map[string]any)
	schema := jsonSchemaBody["schema"].(map[string]any)
	properties := schema["properties"].(map[string]any)
	list := properties["recommendations"].(map[string]any)
	if list["minItems"] != 1 || list["maxItems"] != 2 {
		t.Fatalf("recommendations must have one or two items in the provider schema: %#v", list)
	}
	items := list["items"].(map[string]any)
	if len(items["enum"].([]string)) != 2 {
		t.Fatalf("provider schema must constrain recommendations to verified options: %#v", items)
	}
}

func TestScenarioClaimsRejectImplementedWorkAndUnsupportedBudgetJudgment(t *testing.T) {
	for _, phrase := range []string{
		"Цифровая платформа уже внедрена.",
		"В районе реализуются выбранные проекты.",
		"Почти весь бюджет использован эффективно.",
		"Бюджет говорит о высокой активности.",
		"Критические проблемы устранены.",
	} {
		if validScenarioClaim(phrase) {
			t.Errorf("unsupported scenario claim accepted: %q", phrase)
		}
	}
	if !validScenarioClaim("По расчёту в сценарии выбрана цифровая платформа, а бюджет почти распределён.") {
		t.Fatal("a factual hypothetical claim must remain valid")
	}
}

func TestCompareSendsVerifiedFactsOnly(t *testing.T) {
	requests := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		var body struct {
			Messages []struct {
				Content string `json:"content"`
			} `json:"messages"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body.Messages) != 2 || !strings.Contains(body.Messages[1].Content, "verifiedFacts") || strings.Contains(body.Messages[1].Content, "scoreAfter") {
			t.Fatalf("provider received raw result: %#v", body.Messages)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": `{"summary":"Сценарий А лидирует по итоговому показателю."}`}}}})
	}))
	defer provider.Close()
	c := &Client{baseURL: provider.URL, apiKey: "test", model: "test", http: provider.Client()}
	scenarios := []NamedResult{{Label: "А", Result: testResult()}, {Label: "Б", Result: testResult()}}
	got := c.Compare(context.Background(), scenarios)
	if requests != 1 || got.Source != "live" {
		t.Fatalf("want live comparison; requests=%d result=%#v", requests, got)
	}
}

func testResult() *scoring.Result {
	score := 56.5431
	return &scoring.Result{
		Score: &score, DAvgBefore: 56.8624, DAvgAfter: 58.0776,
		CriticalBefore: 2, CriticalAfter: 0,
		BudgetUsed: 95, BudgetRemaining: 5,
		MinDistrictID: "nura",
		Districts:     []scoring.DistrictResult{{ID: "nura", Name: "Нура"}},
	}
}
