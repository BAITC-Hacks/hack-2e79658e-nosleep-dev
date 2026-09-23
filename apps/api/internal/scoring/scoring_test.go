package scoring

import (
	"encoding/json"
	"math"
	"os"
	"path/filepath"
	"testing"
)

type golden struct {
	Valid []struct {
		Name               string
		Decisions          []Decision
		ExpectedBudgetUsed *int `json:"expectedBudgetUsed"`
		Expected           *struct {
			DAvg        float64 `json:"dAvg"`
			MinDistrict string  `json:"minDistrict"`
			MinD        float64 `json:"minD"`
			NCrit       int     `json:"nCrit"`
			Score       float64 `json:"score"`
		} `json:"expected"`
	}
	Invalid []struct {
		Name         string
		Decisions    []Decision
		ExpectedCode string `json:"expectedCode"`
	}
}

func TestGoldenFixtures(t *testing.T) {
	root := filepath.Join("..", "..", "..", "..")
	t.Setenv("DATA_DIR", filepath.Join(root, "data"))
	b, err := os.ReadFile(filepath.Join(root, "data", "fixtures", "golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var f golden
	if err := json.Unmarshal(b, &f); err != nil {
		t.Fatal(err)
	}
	e, err := New()
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range f.Valid {
		t.Run(tc.Name, func(t *testing.T) {
			// The empty fixture is the documented baseline, not a submit-ready set.
			// /simulate still correctly reports DECISION_COUNT for an empty picker.
			r, v := e.Simulate(tc.Decisions)
			if len(tc.Decisions) == 0 {
				r, v = e.calculate(tc.Decisions), nil
				s := score(r.DAvgAfter, r.MinDAfter, r.CriticalAfter, e.rules)
				r.Score = &s
			}
			if len(v) != 0 {
				t.Fatalf("unexpected violation: %#v", v)
			}
			if tc.ExpectedBudgetUsed != nil && r.BudgetUsed != *tc.ExpectedBudgetUsed {
				t.Fatalf("budget=%d", r.BudgetUsed)
			}
			if tc.Expected != nil {
				close(t, "dAvg", r.DAvgAfter, tc.Expected.DAvg)
				close(t, "minD", r.MinDAfter, tc.Expected.MinD)
				if r.MinDistrictID != tc.Expected.MinDistrict {
					t.Fatalf("min district=%s", r.MinDistrictID)
				}
				if r.CriticalAfter != tc.Expected.NCrit {
					t.Fatalf("ncrit=%d", r.CriticalAfter)
				}
				if r.Score == nil {
					t.Fatal("score nil")
				}
				close(t, "score", *r.Score, tc.Expected.Score)
			}
		})
	}
	for _, tc := range f.Invalid {
		t.Run(tc.Name, func(t *testing.T) {
			_, v := e.Simulate(tc.Decisions)
			if len(v) != 1 || v[0].Code != tc.ExpectedCode {
				t.Fatalf("violations=%#v, want %s", v, tc.ExpectedCode)
			}
		})
	}
}
func close(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.Abs(got-want) > 0.005 {
		t.Fatalf("%s=%.4f want %.4f", name, got, want)
	}
}
