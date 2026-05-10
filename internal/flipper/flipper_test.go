package flipper

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	s.Set("dark-mode", true)

	f, err := s.Get("dark-mode")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Name != "dark-mode" {
		t.Errorf("expected name dark-mode, got %s", f.Name)
	}
	if !f.Enabled {
		t.Error("expected flag to be enabled")
	}
	if f.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, err := s.Get("missing")
	if err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestIsEnabled_True(t *testing.T) {
	s := New()
	s.Set("beta", true)
	if !s.IsEnabled("beta") {
		t.Error("expected IsEnabled to return true")
	}
}

func TestIsEnabled_False(t *testing.T) {
	s := New()
	s.Set("beta", false)
	if s.IsEnabled("beta") {
		t.Error("expected IsEnabled to return false")
	}
}

func TestIsEnabled_Missing(t *testing.T) {
	s := New()
	if s.IsEnabled("nonexistent") {
		t.Error("expected IsEnabled to return false for missing flag")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	s.Set("feature", true)
	s.Set("feature", false)

	f, err := s.Get("feature")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Enabled {
		t.Error("expected flag to be disabled after overwrite")
	}
}

func TestDelete_RemovesFlag(t *testing.T) {
	s := New()
	s.Set("temp", true)
	if err := s.Delete("temp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err := s.Get("temp")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestDelete_Unknown(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error when deleting unknown flag")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	s.Set("a", true)
	s.Set("b", false)

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 flags, got %d", len(all))
	}
}
