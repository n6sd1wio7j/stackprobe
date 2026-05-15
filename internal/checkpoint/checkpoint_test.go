package checkpoint

import (
	"testing"
)

func makeRecord(t *testing.T, s *Store, service, status, note string) {
	t.Helper()
	if err := s.Save(service, status, note); err != nil {
		t.Fatalf("Save: %v", err)
	}
}

func TestSave_And_Get(t *testing.T) {
	s := New()
	makeRecord(t, s, "api", "up", "all good")
	r, err := s.Get("api")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if r.Service != "api" || r.Status != "up" || r.Note != "all good" {
		t.Errorf("unexpected record: %+v", r)
	}
	if r.SavedAt.IsZero() {
		t.Error("SavedAt should be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, err := s.Get("ghost")
	if err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestSave_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Save("", "up", ""); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestSave_EmptyStatusReturnsError(t *testing.T) {
	s := New()
	if err := s.Save("api", "", ""); err == nil {
		t.Error("expected error for empty status")
	}
}

func TestSave_Overwrites(t *testing.T) {
	s := New()
	makeRecord(t, s, "api", "up", "first")
	makeRecord(t, s, "api", "down", "second")
	r, _ := s.Get("api")
	if r.Status != "down" || r.Note != "second" {
		t.Errorf("expected overwritten record, got %+v", r)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	makeRecord(t, s, "api", "up", "")
	if err := s.Delete("api"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := s.Get("api"); err == nil {
		t.Error("expected error after deletion")
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
	makeRecord(t, s, "a", "up", "")
	makeRecord(t, s, "b", "down", "")
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
}
