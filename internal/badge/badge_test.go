package badge

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("api", "status", "up", "green", StyleFlat); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := s.Get("api")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.Status != "up" || b.Color != "green" || b.Label != "status" {
		t.Errorf("unexpected badge: %+v", b)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, err := s.Get("ghost")
	if err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", "status", "up", "green", StyleFlat); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_DefaultsColor(t *testing.T) {
	s := New()
	if err := s.Set("svc", "label", "unknown", "", StyleFlat); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := s.Get("svc")
	if b.Color != "grey" {
		t.Errorf("expected default color grey, got %q", b.Color)
	}
}

func TestSet_DefaultsStyle(t *testing.T) {
	s := New()
	if err := s.Set("svc", "label", "up", "green", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := s.Get("svc")
	if b.Style != StyleFlat {
		t.Errorf("expected default style flat, got %q", b.Style)
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	_ = s.Set("api", "status", "up", "green", StyleFlat)
	_ = s.Set("api", "status", "down", "red", StylePlastic)
	b, _ := s.Get("api")
	if b.Status != "down" || b.Style != StylePlastic {
		t.Errorf("expected overwritten badge, got %+v", b)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("api", "status", "up", "green", StyleFlat)
	if err := s.Delete("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.Get("api"); err == nil {
		t.Fatal("expected error after deletion")
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
	_ = s.Set("a", "l", "up", "green", StyleFlat)
	_ = s.Set("b", "l", "down", "red", StyleFlat)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 badges, got %d", len(all))
	}
}
