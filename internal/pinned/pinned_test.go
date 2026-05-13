package pinned

import (
	"testing"
)

func TestPin_And_Get(t *testing.T) {
	s := New()
	if err := s.Pin("svc-a", true, "manual override"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected pinned entry to exist")
	}
	if v.Service != "svc-a" || !v.Up || v.Reason != "manual override" {
		t.Errorf("unexpected status: %+v", v)
	}
	if v.PinnedAt.IsZero() {
		t.Error("PinnedAt should be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("ghost")
	if ok {
		t.Error("expected no entry for unknown service")
	}
}

func TestPin_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Pin("", true, "reason"); err == nil {
		t.Error("expected error for empty service name")
	}
}

func TestUnpin_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Pin("svc-b", false, "degraded")
	if err := s.Unpin("svc-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsPinned("svc-b") {
		t.Error("expected service to be unpinned")
	}
}

func TestUnpin_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Unpin("nonexistent"); err == nil {
		t.Error("expected error when unpinning unknown service")
	}
}

func TestIsPinned_True(t *testing.T) {
	s := New()
	_ = s.Pin("svc-c", true, "")
	if !s.IsPinned("svc-c") {
		t.Error("expected service to be pinned")
	}
}

func TestIsPinned_False(t *testing.T) {
	s := New()
	if s.IsPinned("svc-d") {
		t.Error("expected service not to be pinned")
	}
}

func TestAll_ReturnsAllPinned(t *testing.T) {
	s := New()
	_ = s.Pin("x", true, "a")
	_ = s.Pin("y", false, "b")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}

func TestPin_Overwrites(t *testing.T) {
	s := New()
	_ = s.Pin("svc-e", true, "first")
	_ = s.Pin("svc-e", false, "second")
	v, _ := s.Get("svc-e")
	if v.Up || v.Reason != "second" {
		t.Errorf("expected overwritten value, got %+v", v)
	}
}
