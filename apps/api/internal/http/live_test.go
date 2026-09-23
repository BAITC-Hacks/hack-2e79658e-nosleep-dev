package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLiveServicesServeCalculatedScenario(t *testing.T) {
	t.Setenv("DEMO_MODE", "true")
	engine, explainer, deck, err := NewLiveServices()
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(engine, explainer, deck, DefaultConfig())
	get := func(path string) *httptest.ResponseRecorder {
		t.Helper()
		response := httptest.NewRecorder()
		router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
		if response.Code != http.StatusOK {
			t.Fatalf("%s returned %d: %s", path, response.Code, response.Body.String())
		}
		return response
	}

	post := func(path, body string) *httptest.ResponseRecorder {
		t.Helper()
		request := httptest.NewRequest(http.MethodPost, path, strings.NewReader(body))
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("%s returned %d: %s", path, response.Code, response.Body.String())
		}
		return response
	}
	var catalog Catalog
	if err := json.Unmarshal(get("/api/v1/catalog").Body.Bytes(), &catalog); err != nil {
		t.Fatal(err)
	}
	if len(catalog.Districts) != 5 || catalog.RequiredDecisions != 5 {
		t.Fatalf("unexpected live catalog: %+v", catalog)
	}
	var optimum OptimumResult
	if err := json.Unmarshal(get("/api/v1/optimum").Body.Bytes(), &optimum); err != nil {
		t.Fatal(err)
	}
	if optimum.BestScore <= 0 || len(optimum.Decisions) != 0 {
		t.Fatalf("unexpected unrevealed optimum: %+v", optimum)
	}

	baseline := post("/api/v1/simulate", `{"decisions":[]}`)
	partial := post("/api/v1/simulate", `{"decisions":[{"initiativeId":"M7","districtId":"nura"}]}`)
	var before, after simulateResponse
	if err := json.Unmarshal(baseline.Body.Bytes(), &before); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(partial.Body.Bytes(), &after); err != nil {
		t.Fatal(err)
	}
	if after.DAvgAfter <= before.DAvgAfter || after.BudgetUsed != 24 || after.Submittable {
		t.Fatalf("partial decision was not calculated: %+v", after.Result)
	}

	valid := `{"decisions":[{"initiativeId":"M7","districtId":"nura"},{"initiativeId":"M8","districtId":"nura"},{"initiativeId":"M10","districtId":"nura"},{"initiativeId":"M12"},{"initiativeId":"M5","districtId":"saryarka"}]}`
	var explanation explainResponse
	if err := json.Unmarshal(post("/api/v1/explain", valid).Body.Bytes(), &explanation); err != nil {
		t.Fatal(err)
	}
	if explanation.Result.Score == nil || explanation.Source != "cached" || explanation.Explanation.Summary == "" {
		t.Fatalf("expected a calculated score and labeled explanation: %+v", explanation)
	}
}
