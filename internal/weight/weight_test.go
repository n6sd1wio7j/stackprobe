package weight

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New(1)
	if err := s.Set("api", 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("api"); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(3)
	if got := s.Get("unknown"); got != 3 {
		t.Fatalf("expected default 3, got %d", got)
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New(1)
	if err := s.Set("", 5); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_ZeroWeightReturnsError(t *testing.T) {
	s := New(1)
	if err := s.Set("api", 0); err == nil {
		t.Fatal("expected error for zero weight")
	}
}

func TestSet_NegativeWeightReturnsError(t *testing.T) {
	s := New(1)
	if err := s.Set("api", -2); err == nil {
		t.Fatal("expected error for negative weight")
	}
}

func TestDelete_RevertsToDefault(t *testing.T) {
	s := New(2)
	_ = s.Set("svc", 10)
	s.Delete("svc")
	if got := s.Get("svc"); got != 2 {
		t.Fatalf("expected default 2 after delete, got %d", got)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New(1)
	_ = s.Set("a", 2)
	_ = s.Set("b", 4)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["a"] != 2 || all["b"] != 4 {
		t.Fatalf("unexpected values: %v", all)
	}
}

func TestDefault_ReturnsConfiguredValue(t *testing.T) {
	s := New(7)
	if s.Default() != 7 {
		t.Fatalf("expected default 7, got %d", s.Default())
	}
}

func TestNew_NegativeDefaultClampsToOne(t *testing.T) {
	s := New(-5)
	if s.Default() != 1 {
		t.Fatalf("expected clamped default 1, got %d", s.Default())
	}
}
