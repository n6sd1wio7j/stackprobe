package mute

import (
	"testing"
	"time"
)

func futureTime(d time.Duration) time.Time {
	return time.Now().Add(d)
}

func TestMute_And_IsMuted(t *testing.T) {
	s := New()
	if err := s.Mute("svc-a", "maintenance", futureTime(time.Hour)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsMuted("svc-a") {
		t.Fatal("expected svc-a to be muted")
	}
}

func TestMute_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Mute("", "reason", futureTime(time.Hour)); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestMute_PastTimeReturnsError(t *testing.T) {
	s := New()
	if err := s.Mute("svc", "reason", time.Now().Add(-time.Minute)); err == nil {
		t.Fatal("expected error for past until time")
	}
}

func TestUnmute_RemovesWindow(t *testing.T) {
	s := New()
	_ = s.Mute("svc-b", "test", futureTime(time.Hour))
	if err := s.Unmute("svc-b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsMuted("svc-b") {
		t.Fatal("expected svc-b to be unmuted")
	}
}

func TestUnmute_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Unmute("no-such-svc"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestIsMuted_UnknownService(t *testing.T) {
	s := New()
	if s.IsMuted("ghost") {
		t.Fatal("expected false for unknown service")
	}
}

func TestAll_ReturnsActiveMutes(t *testing.T) {
	s := New()
	_ = s.Mute("svc-1", "r1", futureTime(time.Hour))
	_ = s.Mute("svc-2", "r2", futureTime(2*time.Hour))
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 mute windows, got %d", len(all))
	}
}

func TestGet_KnownService(t *testing.T) {
	s := New()
	_ = s.Mute("svc-x", "deploy", futureTime(time.Hour))
	w, ok := s.Get("svc-x")
	if !ok {
		t.Fatal("expected window to be found")
	}
	if w.Reason != "deploy" {
		t.Fatalf("expected reason 'deploy', got %q", w.Reason)
	}
}
