package dependency

import (
	"testing"
)

func TestAdd_And_Deps(t *testing.T) {
	s := New()
	if err := s.Add("api", "db"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	deps := s.Deps("api")
	if len(deps) != 1 || deps[0] != "db" {
		t.Errorf("expected [db], got %v", deps)
	}
}

func TestDeps_Unknown(t *testing.T) {
	s := New()
	if deps := s.Deps("unknown"); len(deps) != 0 {
		t.Errorf("expected empty slice, got %v", deps)
	}
}

func TestAdd_DetectsCycle(t *testing.T) {
	s := New()
	_ = s.Add("a", "b")
	_ = s.Add("b", "c")
	if err := s.Add("c", "a"); err != ErrCycle {
		t.Errorf("expected ErrCycle, got %v", err)
	}
}

func TestAdd_DirectCycle(t *testing.T) {
	s := New()
	_ = s.Add("x", "y")
	if err := s.Add("y", "x"); err != ErrCycle {
		t.Errorf("expected ErrCycle for direct cycle, got %v", err)
	}
}

func TestAffected_DirectDependent(t *testing.T) {
	s := New()
	_ = s.Add("api", "db")
	affected := s.Affected("db")
	if len(affected) != 1 || affected[0] != "api" {
		t.Errorf("expected [api], got %v", affected)
	}
}

func TestAffected_Transitive(t *testing.T) {
	s := New()
	_ = s.Add("frontend", "api")
	_ = s.Add("api", "db")
	affected := s.Affected("db")
	if len(affected) != 2 {
		t.Errorf("expected 2 affected services, got %v", affected)
	}
	seen := map[string]bool{}
	for _, a := range affected {
		seen[a] = true
	}
	if !seen["api"] || !seen["frontend"] {
		t.Errorf("expected api and frontend in affected, got %v", affected)
	}
}

func TestAffected_None(t *testing.T) {
	s := New()
	_ = s.Add("api", "db")
	if affected := s.Affected("cache"); len(affected) != 0 {
		t.Errorf("expected no affected services, got %v", affected)
	}
}

func TestRemove_ClearsDeps(t *testing.T) {
	s := New()
	_ = s.Add("api", "db")
	s.Remove("api")
	if deps := s.Deps("api"); len(deps) != 0 {
		t.Errorf("expected empty deps after remove, got %v", deps)
	}
}
