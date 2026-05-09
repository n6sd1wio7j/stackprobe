package banner

import (
	"testing"
	"time"
)

func TestAdd_And_Active(t *testing.T) {
	s := New()
	id := s.Add("Scheduled maintenance tonight", "warning", time.Time{})
	if id == "" {
		t.Fatal("expected non-empty id")
	}
	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active banner, got %d", len(active))
	}
	if active[0].Message != "Scheduled maintenance tonight" {
		t.Errorf("unexpected message: %s", active[0].Message)
	}
}

func TestAdd_ExpiredNotReturned(t *testing.T) {
	s := New()
	past := time.Now().Add(-1 * time.Hour)
	s.Add("Old message", "info", past)
	active := s.Active()
	if len(active) != 0 {
		t.Fatalf("expected 0 active banners, got %d", len(active))
	}
}

func TestRemove_KnownID(t *testing.T) {
	s := New()
	id := s.Add("hello", "info", time.Time{})
	ok := s.Remove(id)
	if !ok {
		t.Fatal("expected Remove to return true")
	}
	if len(s.Active()) != 0 {
		t.Fatal("expected no active banners after removal")
	}
}

func TestRemove_UnknownID(t *testing.T) {
	s := New()
	ok := s.Remove("banner-999")
	if ok {
		t.Fatal("expected Remove to return false for unknown id")
	}
}

func TestActive_MultipleWithMixedExpiry(t *testing.T) {
	s := New()
	s.Add("permanent", "info", time.Time{})
	s.Add("future", "warning", time.Now().Add(time.Hour))
	s.Add("expired", "critical", time.Now().Add(-time.Minute))

	active := s.Active()
	if len(active) != 2 {
		t.Fatalf("expected 2 active banners, got %d", len(active))
	}
}

func TestAdd_AssignsUniqueIDs(t *testing.T) {
	s := New()
	id1 := s.Add("first", "info", time.Time{})
	id2 := s.Add("second", "info", time.Time{})
	if id1 == id2 {
		t.Errorf("expected unique IDs, got %s and %s", id1, id2)
	}
}
