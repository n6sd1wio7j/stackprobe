package grouping

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("backend", []string{"api", "db"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	svcs, ok := s.Get("backend")
	if !ok {
		t.Fatal("expected group to exist")
	}
	if len(svcs) != 2 || svcs[0] != "api" || svcs[1] != "db" {
		t.Fatalf("unexpected services: %v", svcs)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected false for unknown group")
	}
}

func TestSet_EmptyNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", []string{"svc"}); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestSet_EmptyServicesReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("g", []string{}); err == nil {
		t.Fatal("expected error for empty services")
	}
}

func TestSet_EmptyServiceNameReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("g", []string{"ok", ""}); err == nil {
		t.Fatal("expected error for blank service name")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	_ = s.Set("g", []string{"a"})
	_ = s.Set("g", []string{"b", "c"})
	svcs, _ := s.Get("g")
	if len(svcs) != 2 {
		t.Fatalf("expected 2 services, got %d", len(svcs))
	}
}

func TestDelete_RemovesGroup(t *testing.T) {
	s := New()
	_ = s.Set("g", []string{"a"})
	if err := s.Delete("g"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("g")
	if ok {
		t.Fatal("expected group to be deleted")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("nope"); err == nil {
		t.Fatal("expected error for unknown group")
	}
}

func TestMember_True(t *testing.T) {
	s := New()
	_ = s.Set("g", []string{"alpha", "beta"})
	if !s.Member("g", "alpha") {
		t.Fatal("expected alpha to be a member")
	}
}

func TestMember_False(t *testing.T) {
	s := New()
	_ = s.Set("g", []string{"alpha"})
	if s.Member("g", "gamma") {
		t.Fatal("expected gamma not to be a member")
	}
}

func TestAll_ReturnsAllGroups(t *testing.T) {
	s := New()
	_ = s.Set("g1", []string{"a"})
	_ = s.Set("g2", []string{"b", "c"})
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(all))
	}
}
