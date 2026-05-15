package suppress

import (
	"testing"
	"time"
)

func futureTime() time.Time { return time.Now().Add(time.Hour) }

func TestSuppress_And_IsSuppressed(t *testing.T) {
	s := New()
	if err := s.Suppress("svc", "maintenance", futureTime()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsSuppressed("svc") {
		t.Fatal("expected service to be suppressed")
	}
}

func TestSuppress_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Suppress("", "reason", futureTime()); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSuppress_EmptyReasonReturnsError(t *testing.T) {
	s := New()
	if err := s.Suppress("svc", "", futureTime()); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestSuppress_PastTimeReturnsError(t *testing.T) {
	s := New()
	past := time.Now().Add(-time.Minute)
	if err := s.Suppress("svc", "reason", past); err == nil {
		t.Fatal("expected error for past expiry")
	}
}

func TestUnsuppress_RemovesWindow(t *testing.T) {
	s := New()
	_ = s.Suppress("svc", "reason", futureTime())
	if err := s.Unsuppress("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsSuppressed("svc") {
		t.Fatal("expected service to be unsuppressed")
	}
}

func TestUnsuppress_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Unsuppress("missing"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestIsSuppressed_Unknown(t *testing.T) {
	s := New()
	if s.IsSuppressed("ghost") {
		t.Fatal("unknown service should not be suppressed")
	}
}

func TestIsSuppressed_Expired(t *testing.T) {
	s := New()
	// Insert a rule that is already expired by manipulating directly.
	s.rules["svc"] = Rule{
		Service:   "svc",
		Reason:    "old",
		ExpiresAt: time.Now().Add(-time.Minute),
		CreatedAt: time.Now().Add(-2 * time.Minute),
	}
	if s.IsSuppressed("svc") {
		t.Fatal("expired rule should not count as suppressed")
	}
}

func TestAll_ReturnsAllRules(t *testing.T) {
	s := New()
	_ = s.Suppress("a", "r", futureTime())
	_ = s.Suppress("b", "r", futureTime())
	if got := len(s.All()); got != 2 {
		t.Fatalf("expected 2 rules, got %d", got)
	}
}
