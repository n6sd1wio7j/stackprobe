package endpointmeta

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	m := Meta{Description: "main API", Owner: "platform", Tier: "1"}
	if err := s.Set("api", m); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("api")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Service != "api" {
		t.Errorf("service = %q, want %q", got.Service, "api")
	}
	if got.Tier != "1" {
		t.Errorf("tier = %q, want %q", got.Tier, "1")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", Meta{}); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	_ = s.Set("svc", Meta{Tier: "2"})
	_ = s.Set("svc", Meta{Tier: "1"})
	got, _ := s.Get("svc")
	if got.Tier != "1" {
		t.Errorf("tier = %q, want %q", got.Tier, "1")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("svc", Meta{Description: "temp"})
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected entry to be removed")
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
	_ = s.Set("a", Meta{Tier: "1"})
	_ = s.Set("b", Meta{Tier: "2"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("len = %d, want 2", len(all))
	}
}
