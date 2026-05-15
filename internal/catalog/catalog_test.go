package catalog

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	s := New()
	err := s.Register("svc-a", "Service A", "team-x", "http://svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Service != "svc-a" || e.URL != "http://svc-a" {
		t.Errorf("unexpected entry: %+v", e)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestRegister_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("", "desc", "owner", "http://x"); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestRegister_EmptyURLReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("svc", "desc", "owner", ""); err == nil {
		t.Fatal("expected error for empty url")
	}
}

func TestRegister_Overwrites(t *testing.T) {
	s := New()
	_ = s.Register("svc", "old", "team", "http://old")
	_ = s.Register("svc", "new", "team", "http://new")
	e, _ := s.Get("svc")
	if e.Description != "new" || e.URL != "http://new" {
		t.Errorf("expected updated entry, got %+v", e)
	}
	if e.Registered.IsZero() {
		t.Error("registered time should not be zero")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Register("svc", "desc", "owner", "http://svc")
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("svc")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := New()
	_ = s.Register("a", "", "", "http://a")
	_ = s.Register("b", "", "", "http://b")
	if len(s.All()) != 2 {
		t.Errorf("expected 2 entries, got %d", len(s.All()))
	}
}
