package runbook

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	e, err := s.Set("svc-a", "Restart guide", "https://wiki.example.com/svc-a", "Check logs first")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "svc-a" {
		t.Errorf("expected service svc-a, got %s", e.Service)
	}
	got, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Title != "Restart guide" {
		t.Errorf("expected title 'Restart guide', got %s", got.Title)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	_, err := s.Set("", "title", "", "")
	if err == nil {
		t.Error("expected error for empty service")
	}
}

func TestSet_EmptyTitleReturnsError(t *testing.T) {
	s := New()
	_, err := s.Set("svc-a", "", "", "")
	if err == nil {
		t.Error("expected error for empty title")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	s.Set("svc-a", "Old title", "", "")
	s.Set("svc-a", "New title", "https://new", "updated")
	got, _ := s.Get("svc-a")
	if got.Title != "New title" {
		t.Errorf("expected 'New title', got %s", got.Title)
	}
	if got.URL != "https://new" {
		t.Errorf("expected url 'https://new', got %s", got.URL)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	s.Set("svc-a", "Guide", "", "")
	if err := s.Delete("svc-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("svc-a")
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestAll_ReturnsAllEntries(t *testing.T) {
	s := New()
	s.Set("svc-a", "Guide A", "", "")
	s.Set("svc-b", "Guide B", "", "")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
