package httpapi

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()
	gin.SetMode(gin.TestMode)
	services, err := NewFixtureServices("")
	if err != nil {
		t.Fatal(err)
	}
	config := DefaultConfig()
	return NewRouter(services, services, services, config)
}

func TestHealthAndCORS(t *testing.T) {
	handler := testRouter(t)
	request := httptest.NewRequest(http.MethodGet, "/health", nil)
	request.Header.Set("Origin", "http://localhost:3000")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || recorder.Body.String() != "{\"status\":\"ok\"}" {
		t.Fatalf("health response = %d %q", recorder.Code, recorder.Body.String())
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Fatalf("Access-Control-Allow-Origin = %q", got)
	}
}

func TestSimulateHappyPath(t *testing.T) {
	server := testRouter(t)

	body := []byte(`{"decisions":[{"initiativeId":"M7","districtId":"nura"},{"initiativeId":"M8","districtId":"nura"},{"initiativeId":"M10","districtId":"nura"},{"initiativeId":"M12"},{"initiativeId":"M5","districtId":"saryarka"}]}`)
	response := postJSON(t, server, "/api/v1/simulate", body)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	var payload struct {
		Submittable bool        `json:"submittable"`
		Violations  []Violation `json:"violations"`
		BudgetUsed  int         `json:"budgetUsed"`
		Score       *float64    `json:"score"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !payload.Submittable || len(payload.Violations) != 0 {
		t.Fatalf("scenario did not submit: %+v", payload)
	}
	if payload.BudgetUsed != 95 || payload.Score == nil || math.Abs(*payload.Score-56.5431) > 0.0001 {
		t.Fatalf("unexpected fixture result: budget=%d score=%v", payload.BudgetUsed, payload.Score)
	}
}

func TestExplainBadDecisionCountUsesContractErrorShape(t *testing.T) {
	server := testRouter(t)

	body := []byte(`{"decisions":[{"initiativeId":"M7","districtId":"nura"},{"initiativeId":"M12"}]}`)
	response := postJSON(t, server, "/api/v1/explain", body)
	defer response.Body.Close()
	if response.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want 422", response.StatusCode)
	}
	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
			Details struct {
				Violations []Violation `json:"violations"`
			} `json:"details"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Error.Code != "DECISION_COUNT" || payload.Error.Message == "" || len(payload.Error.Details.Violations) == 0 {
		t.Fatalf("unexpected error payload: %+v", payload.Error)
	}
}

func TestMalformedRequestTypeReturnsContractError(t *testing.T) {
	server := testRouter(t)

	response := postJSON(t, server, "/api/v1/simulate", []byte(`{"decisions":[{"initiativeId":7}]}`))
	defer response.Body.Close()
	if response.StatusCode != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", response.StatusCode)
	}
	var payload struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Error.Code != "INVALID_JSON" || payload.Error.Message == "" {
		t.Fatalf("unexpected error payload: %+v", payload.Error)
	}
}

func TestSimulateGoldenFixtureRoundTrip(t *testing.T) {
	server := testRouter(t)

	repoRoot := findRepoRoot(t)
	goldenBytes, err := os.ReadFile(filepath.Join(repoRoot, "data", "fixtures", "golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var golden struct {
		Valid []struct {
			Name               string     `json:"name"`
			Decisions          []Decision `json:"decisions"`
			ExpectedBudgetUsed int        `json:"expectedBudgetUsed"`
			Expected           struct {
				Score float64 `json:"score"`
			} `json:"expected"`
		} `json:"valid"`
	}
	if err := json.Unmarshal(goldenBytes, &golden); err != nil {
		t.Fatal(err)
	}
	var fixture struct {
		Name               string
		Decisions          []Decision
		ExpectedBudgetUsed int
		ExpectedScore      float64
	}
	for _, row := range golden.Valid {
		if row.Name == "brief_example_set" {
			fixture.Name = row.Name
			fixture.Decisions = row.Decisions
			fixture.ExpectedBudgetUsed = row.ExpectedBudgetUsed
			fixture.ExpectedScore = row.Expected.Score
			break
		}
	}
	if fixture.Name == "" {
		t.Fatal("brief_example_set missing from golden fixture")
	}
	body, err := json.Marshal(struct {
		Decisions []Decision `json:"decisions"`
	}{Decisions: fixture.Decisions})
	if err != nil {
		t.Fatal(err)
	}
	response := postJSON(t, server, "/api/v1/simulate", body)
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	var payload struct {
		Submittable bool     `json:"submittable"`
		BudgetUsed  int      `json:"budgetUsed"`
		Score       *float64 `json:"score"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if !payload.Submittable || payload.BudgetUsed != fixture.ExpectedBudgetUsed || payload.Score == nil || math.Abs(*payload.Score-fixture.ExpectedScore) > 0.0001 {
		t.Fatalf("golden round-trip mismatch: got %+v, want budget=%d score=%.4f", payload, fixture.ExpectedBudgetUsed, fixture.ExpectedScore)
	}
}

func postJSON(t *testing.T, handler http.Handler, path string, body []byte) *http.Response {
	t.Helper()
	request := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	return recorder.Result()
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(current, "data", "fixtures", "golden.json")); err == nil {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			t.Fatal("could not find repository root")
		}
		current = parent
	}
}
