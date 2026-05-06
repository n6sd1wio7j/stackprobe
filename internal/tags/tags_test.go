package tags_test

import (
	"sort"
	"testing"

	"github.com/stackprobe/internal/tags"
)

func TestSet_And_Get(t *testing.T) {
	s := tags.New()
	s.Set("api", []string{"production", "critical"})

	got := s.Get("api")
	if len(got) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(got))
	}
}

func TestGet_Unknown(t *testing.T) {
	s := tags.New()
	if got := s.Get("nonexistent"); got != nil {
		t.Fatalf("expected nil for unknown service, got %v", got)
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := tags.New()
	s.Set("svc", []string{"v1"})
	s.Set("svc", []string{"v2", "v3"})

	got := s.Get("svc")
	if len(got) != 2 || got[0] != "v2" {
		t.Fatalf("expected overwritten tags, got %v", got)
	}
}

func TestFilter_MatchesAll(t *testing.T) {
	s := tags.New()
	s.Set("api", []string{"production", "critical"})
	s.Set("worker", []string{"production"})
	s.Set("db", []string{"critical", "internal"})

	matches := s.Filter("production", "critical")
	if len(matches) != 1 || matches[0] != "api" {
		t.Fatalf("expected [api], got %v", matches)
	}
}

func TestFilter_SingleTag(t *testing.T) {
	s := tags.New()
	s.Set("api", []string{"production", "critical"})
	s.Set("worker", []string{"production"})
	s.Set("db", []string{"internal"})

	matches := s.Filter("production")
	sort.Strings(matches)
	if len(matches) != 2 || matches[0] != "api" || matches[1] != "worker" {
		t.Fatalf("expected [api worker], got %v", matches)
	}
}

func TestFilter_NoMatch(t *testing.T) {
	s := tags.New()
	s.Set("api", []string{"production"})

	matches := s.Filter("staging")
	if len(matches) != 0 {
		t.Fatalf("expected no matches, got %v", matches)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := tags.New()
	s.Set("api", []string{"production"})
	s.Set("db", []string{"internal"})

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	// Mutate returned map — should not affect store
	all["api"] = []string{"mutated"}
	if got := s.Get("api"); len(got) != 1 || got[0] != "production" {
		t.Fatal("store was mutated through returned map")
	}
}
