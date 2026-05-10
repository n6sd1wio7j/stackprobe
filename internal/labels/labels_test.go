package labels

import (
	"sort"
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Set("api", "team", "platform")

	got, err := s.Get("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["env"] != "prod" || got["team"] != "platform" {
		t.Errorf("unexpected labels: %v", got)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, err := s.Get("ghost")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	s.Set("svc", "env", "staging")
	s.Set("svc", "env", "prod")

	got, _ := s.Get("svc")
	if got["env"] != "prod" {
		t.Errorf("expected prod, got %s", got["env"])
	}
}

func TestDelete_RemovesLabel(t *testing.T) {
	s := New()
	s.Set("svc", "env", "prod")
	s.Delete("svc", "env")

	got, _ := s.Get("svc")
	if _, ok := got["env"]; ok {
		t.Error("label should have been deleted")
	}
}

func TestFilter_MatchesAll(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Set("api", "team", "platform")
	s.Set("worker", "env", "prod")
	s.Set("worker", "team", "data")
	s.Set("db", "env", "staging")

	results := s.Filter(map[string]string{"env": "prod"})
	sort.Strings(results)
	if len(results) != 2 || results[0] != "api" || results[1] != "worker" {
		t.Errorf("unexpected filter results: %v", results)
	}
}

func TestFilter_MultipleSelectors(t *testing.T) {
	s := New()
	s.Set("api", "env", "prod")
	s.Set("api", "team", "platform")
	s.Set("worker", "env", "prod")
	s.Set("worker", "team", "data")

	results := s.Filter(map[string]string{"env": "prod", "team": "platform"})
	if len(results) != 1 || results[0] != "api" {
		t.Errorf("unexpected filter results: %v", results)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	s.Set("a", "k", "v")
	s.Set("b", "k", "v")

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 services, got %d", len(all))
	}
}
