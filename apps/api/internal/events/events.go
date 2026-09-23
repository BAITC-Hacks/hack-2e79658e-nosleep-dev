// Package events exposes a small deterministic-or-random crisis-event deck.
package events

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Event struct {
	ID          string             `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	DistrictID  string             `json:"districtId"`
	Shocks      map[string]float64 `json:"shocks"`
	BudgetDelta int                `json:"budgetDelta"`
}
type Deck interface{ Draw(seed *int64) Event }
type deck struct{ events []Event }

func New() (*deck, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}
	raw, err := os.ReadFile(filepath.Join(dir, "events.json"))
	if err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}
	var events []Event
	if err := json.Unmarshal(raw, &events); err != nil {
		return nil, fmt.Errorf("decode events: %w", err)
	}
	if len(events) == 0 {
		return nil, fmt.Errorf("event deck is empty")
	}
	return &deck{events: events}, nil
}
func MustNew() *deck {
	d, err := New()
	if err != nil {
		panic(err)
	}
	return d
}
func (d *deck) Draw(seed *int64) Event {
	source := time.Now().UnixNano()
	if seed != nil {
		source = *seed
	}
	return d.events[rand.New(rand.NewSource(source)).Intn(len(d.events))]
}
func dataDir() (string, error) {
	if dir := os.Getenv("DATA_DIR"); dir != "" {
		return dir, nil
	}
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		candidate := filepath.Join(wd, "data", "events.json")
		if _, err := os.Stat(candidate); err == nil {
			return filepath.Dir(candidate), nil
		}
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	return "", fmt.Errorf("data directory not found; set DATA_DIR")
}
