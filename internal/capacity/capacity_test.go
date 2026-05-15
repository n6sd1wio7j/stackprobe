package capacity

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("api", 100, 40); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r, ok := s.Get("api")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if r.Limit != 100 || r.Current != 40 {
		t.Fatalf("unexpected values: %+v", r)
	}
	if r.Percent != 40.0 {
		t.Fatalf("expected percent 40.0, got %v", r.Percent)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected no record for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", 10, 5); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_ZeroLimitReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("api", 0, 5); err == nil {
		t.Fatal("expected error for zero limit")
	}
}

func TestSet_NegativeCurrentReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("api", 100, -1); err == nil {
		t.Fatal("expected error for negative current")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("api", 100, 50)
	if err := s.Delete("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("api")
	if ok {
		t.Fatal("expected record to be removed")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Set("a", 10, 1)
	_ = s.Set("b", 20, 5)
	records := s.All()
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
}

func TestIsSaturated_True(t *testing.T) {
	s := New()
	_ = s.Set("api", 10, 10)
	if !s.IsSaturated("api") {
		t.Fatal("expected service to be saturated")
	}
}

func TestIsSaturated_False(t *testing.T) {
	s := New()
	_ = s.Set("api", 10, 5)
	if s.IsSaturated("api") {
		t.Fatal("expected service to not be saturated")
	}
}

func TestIsSaturated_Unknown(t *testing.T) {
	s := New()
	if s.IsSaturated("ghost") {
		t.Fatal("expected false for unknown service")
	}
}
