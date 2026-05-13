package drain

import (
	"testing"
)

func TestEnable_And_IsDraining(t *testing.T) {
	s := New()
	if err := s.Enable("svc-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsDraining("svc-a") {
		t.Fatal("expected svc-a to be draining")
	}
}

func TestEnable_AlreadyDraining(t *testing.T) {
	s := New()
	_ = s.Enable("svc-a")
	if err := s.Enable("svc-a"); err != ErrAlreadyDraining {
		t.Fatalf("expected ErrAlreadyDraining, got %v", err)
	}
}

func TestEnable_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Enable(""); err == nil {
		t.Fatal("expected error for empty service name")
	}
}

func TestDisable_ClearsDrainFlag(t *testing.T) {
	s := New()
	_ = s.Enable("svc-b")
	if err := s.Disable("svc-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsDraining("svc-b") {
		t.Fatal("expected svc-b to no longer be draining")
	}
}

func TestDisable_UnknownService(t *testing.T) {
	s := New()
	if err := s.Disable("ghost"); err != ErrUnknownService {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestGet_ReturnsState(t *testing.T) {
	s := New()
	_ = s.Enable("svc-c")
	st, err := s.Get("svc-c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !st.Draining {
		t.Error("expected Draining to be true")
	}
	if st.StartedAt.IsZero() {
		t.Error("expected StartedAt to be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	if _, err := s.Get("nope"); err != ErrUnknownService {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Enable("svc-x")
	_ = s.Enable("svc-y")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestIsDraining_UnknownReturnsFalse(t *testing.T) {
	s := New()
	if s.IsDraining("unknown") {
		t.Fatal("expected false for unknown service")
	}
}
