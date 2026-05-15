package endpoint

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	s := New()
	if err := s.Register("svc", "http://svc/health", "main service"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	m, err := s.Get("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.URL != "http://svc/health" {
		t.Errorf("expected url %q, got %q", "http://svc/health", m.URL)
	}
	if m.Description != "main service" {
		t.Errorf("expected description %q, got %q", "main service", m.Description)
	}
	if m.RegisteredAt.IsZero() {
		t.Error("expected RegisteredAt to be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, err := s.Get("ghost")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestRegister_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("", "http://x", ""); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestRegister_EmptyURLReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("svc", "", ""); err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestRegister_Overwrites(t *testing.T) {
	s := New()
	_ = s.Register("svc", "http://old", "old")
	_ = s.Register("svc", "http://new", "new")
	m, _ := s.Get("svc")
	if m.URL != "http://new" {
		t.Errorf("expected overwritten url, got %q", m.URL)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Register("svc", "http://svc", "")
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := s.Get("svc")
	if err == nil {
		t.Fatal("expected error after delete")
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
	_ = s.Register("a", "http://a", "")
	_ = s.Register("b", "http://b", "")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
