package envelope

import (
	"testing"
	"time"
)

func TestWrap_And_Get(t *testing.T) {
	s := New(30)
	if err := s.Wrap("svc-a", "up", "all good"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected envelope to exist")
	}
	if e.Status != "up" || e.Message != "all good" {
		t.Errorf("unexpected envelope: %+v", e)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(30)
	_, ok := s.Get("missing")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestWrap_EmptyServiceReturnsError(t *testing.T) {
	s := New(30)
	if err := s.Wrap("", "up", ""); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestWrap_EmptyStatusReturnsError(t *testing.T) {
	s := New(30)
	if err := s.Wrap("svc", "", ""); err == nil {
		t.Error("expected error for empty status")
	}
}

func TestWrap_Overwrites(t *testing.T) {
	s := New(30)
	_ = s.Wrap("svc", "up", "first")
	_ = s.Wrap("svc", "down", "second")
	e, _ := s.Get("svc")
	if e.Status != "down" || e.Message != "second" {
		t.Errorf("expected overwritten envelope, got %+v", e)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New(30)
	_ = s.Wrap("svc", "up", "")
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("svc")
	if ok {
		t.Error("expected envelope to be deleted")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New(30)
	if err := s.Delete("ghost"); err == nil {
		t.Error("expected error for unknown service")
	}
}

func TestAll_ReturnsAllEnvelopes(t *testing.T) {
	s := New(30)
	_ = s.Wrap("a", "up", "")
	_ = s.Wrap("b", "down", "")
	if got := len(s.All()); got != 2 {
		t.Errorf("expected 2 envelopes, got %d", got)
	}
}

func TestIsExpired_Fresh(t *testing.T) {
	s := New(60)
	_ = s.Wrap("svc", "up", "")
	if s.IsExpired("svc") {
		t.Error("expected fresh envelope to not be expired")
	}
}

func TestIsExpired_Stale(t *testing.T) {
	s := New(1)
	_ = s.Wrap("svc", "up", "")
	// Manually back-date the wrapped time.
	s.mu.Lock()
	e := s.items["svc"]
	e.WrappedAt = time.Now().UTC().Add(-5 * time.Second)
	s.items["svc"] = e
	s.mu.Unlock()
	if !s.IsExpired("svc") {
		t.Error("expected stale envelope to be expired")
	}
}

func TestIsExpired_Unknown(t *testing.T) {
	s := New(30)
	if !s.IsExpired("ghost") {
		t.Error("expected unknown service to be treated as expired")
	}
}
