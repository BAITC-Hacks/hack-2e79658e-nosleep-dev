package events

import "testing"

func TestSeededDrawIsDeterministic(t *testing.T) {
	t.Setenv("DATA_DIR", "../../../../data")
	d, err := New()
	if err != nil {
		t.Fatal(err)
	}
	seed := int64(42)
	a, b := d.Draw(&seed), d.Draw(&seed)
	if a.ID != b.ID {
		t.Fatalf("same seed gave %q then %q", a.ID, b.ID)
	}
	if a.ID == "" || a.DistrictID == "" || len(a.Shocks) == 0 {
		t.Fatalf("invalid event: %#v", a)
	}
}
