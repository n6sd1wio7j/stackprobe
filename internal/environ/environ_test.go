package environ_test

import (
	"testing"

	"github.com/yourusername/stackprobe/internal/environ"
)

func TestSet_And_Get(t *testing.T) {
	s := environ.New()
	if err := s.Set("api", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, ok := s.Get("api")
	if !ok || env != "production" {
		t.Fatalf("expected production, got %q ok=%v", env, ok)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := environ.New()
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected not found for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := environ.New()
	if err := s.Set("", "staging"); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_EmptyEnvReturnsError(t *testing.T) {
	s := environ.New()
	if err := s.Set("api", ""); err == nil {
		t.Fatal("expected error for empty environment")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := environ.New()
	_ = s.Set("api", "staging")
	_ = s.Set("api", "production")
	env, _ := s.Get("api")
	if env != "production" {
		t.Fatalf("expected production after overwrite, got %q", env)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := environ.New()
	_ = s.Set("api", "staging")
	s.Delete("api")
	_, ok := s.Get("api")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestFilter_MatchesEnvironment(t *testing.T) {
	s := environ.New()
	_ = s.Set("api", "production")
	_ = s.Set("worker", "staging")
	_ = s.Set("db", "production")
	results := s.Filter("production")
	if len(results) != 2 {
		t.Fatalf("expected 2 production services, got %d", len(results))
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := environ.New()
	_ = s.Set("api", "production")
	_ = s.Set("worker", "staging")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["api"] != "production" || all["worker"] != "staging" {
		t.Fatalf("unexpected snapshot: %v", all)
	}
}
