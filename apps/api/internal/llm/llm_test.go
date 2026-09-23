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
		answer := `{"summary":"Показатель города вырос на 3.78%. Во двух районах остались проблемы.","strengths":["Показатель вырос"],"risks":["Риск"],"recommendations":["Продолжайте работу"]}`
		if requests == 2 {
			answer = `{"summary":"Городской показатель улучшился.","strengths":["Критических показателей не осталось."],"risks":["Бюджет почти исчерпан."],"recommendations":["Уделите внимание самому слабому району."]}`
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"choices": []any{map[string]any{"message": map[string]string{"content": answer}}}})
	}))
	defer provider.Close()
	c := &Client{baseURL: provider.URL, apiKey: "test", model: "test", http: provider.Client()}
	got := c.Explain(context.Background(), nil, testResult())
	if requests != 2 || got.Source != "live" || strings.Contains(got.Summary, "%") {
		t.Fatalf("want grounded live result after retry; requests=%d result=%#v", requests, got)
	}
}

func TestExplainFallsBackAfterUngroundedRetries(t *testing.T) {
	requests := 0
	provider := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		_, _ = w.Write([]byte(`{"choices":[{"message":{"content":"{\"summary\":\"Рост на 3.78%\",\"strengths\":[\"Хорошо\"],\"risks\":[\"Риск\"],\"recommendations\":[\"Продолжайте\"]}"}}]}`))
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

func TestValidPlainTextAllowsVagueQuantity(t *testing.T) {
	if !validPlainText("Для района выбраны несколько инициатив.") {
		t.Fatal("a vague count of selected initiatives should not force a cached fallback")
	}
	if validPlainText("Выбраны пять инициатив.") {
		t.Fatal("exact counts remain excluded from AI prose")
	}
}

func TestExplanationSchemaRequiresNonemptyLists(t *testing.T) {
	jsonSchemaBody := jsonSchema(&Explanation{})["json_schema"].(map[string]any)
	schema := jsonSchemaBody["schema"].(map[string]any)
	properties := schema["properties"].(map[string]any)
	for _, field := range []string{"strengths", "risks", "recommendations"} {
		list := properties[field].(map[string]any)
		if list["minItems"] != 1 {
			t.Fatalf("%s must have at least one item in the provider schema", field)
		}
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
