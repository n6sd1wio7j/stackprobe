package passthrough

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("svc-a", ModePassthrough); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("svc-a"); got != ModePassthrough {
		t.Errorf("expected ModePassthrough, got %v", got)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	if got := s.Get("unknown"); got != ModeNormal {
		t.Errorf("expected ModeNormal for unknown service, got %v", got)
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", ModePassthrough); err == nil {
		t.Error("expected error for empty service name, got nil")
	}
}

func TestIsPassthrough_True(t *testing.T) {
	s := New()
	_ = s.Set("svc-b", ModePassthrough)
	if !s.IsPassthrough("svc-b") {
		t.Error("expected IsPassthrough to return true")
	}
}

func TestIsPassthrough_False(t *testing.T) {
	s := New()
	_ = s.Set("svc-b", ModeNormal)
	if s.IsPassthrough("svc-b") {
		t.Error("expected IsPassthrough to return false")
	}
}

func TestDelete_RevertsToNormal(t *testing.T) {
	s := New()
	_ = s.Set("svc-c", ModePassthrough)
	s.Delete("svc-c")
	if got := s.Get("svc-c"); got != ModeNormal {
		t.Errorf("expected ModeNormal after delete, got %v", got)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Set("svc-x", ModePassthrough)
	_ = s.Set("svc-y", ModeNormal)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["svc-x"] != ModePassthrough {
		t.Errorf("svc-x: expected ModePassthrough, got %v", all["svc-x"])
	}
	if all["svc-y"] != ModeNormal {
		t.Errorf("svc-y: expected ModeNormal, got %v", all["svc-y"])
	}
}

func TestAll_MutationDoesNotAffectStore(t *testing.T) {
	s := New()
	_ = s.Set("svc-z", ModePassthrough)
	all := s.All()
	all["svc-z"] = ModeNormal // mutate snapshot
	if s.Get("svc-z") != ModePassthrough {
		t.Error("store was unexpectedly mutated via snapshot")
	}
}
