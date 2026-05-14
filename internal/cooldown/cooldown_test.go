package cooldown

import (
	"testing"
	"time"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	s := New(time.Minute)
	ok, err := s.Allow("svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	s := New(time.Minute)
	s.Allow("svc-a") //nolint
	ok, err := s.Allow("svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AfterCooldownPermitted(t *testing.T) {
	s := New(10 * time.Millisecond)
	s.Allow("svc-a") //nolint
	time.Sleep(20 * time.Millisecond)
	ok, err := s.Allow("svc-a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected call after cooldown expiry to be allowed")
	}
}

func TestAllow_IndependentServices(t *testing.T) {
	s := New(time.Minute)
	s.Allow("svc-a") //nolint
	ok, err := s.Allow("svc-b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("svc-b should be independent from svc-a cooldown")
	}
}

func TestReset_AllowsImmediateRetrigger(t *testing.T) {
	s := New(time.Minute)
	s.Allow("svc-a") //nolint
	if err := s.Reset("svc-a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ok, _ := s.Allow("svc-a")
	if !ok {
		t.Fatal("expected allow after reset")
	}
}

func TestSet_CustomDuration(t *testing.T) {
	s := New(time.Minute)
	if err := s.Set("svc-a", 10*time.Millisecond); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	s.Allow("svc-a") //nolint
	time.Sleep(20 * time.Millisecond)
	ok, _ := s.Allow("svc-a")
	if !ok {
		t.Fatal("expected allow after custom cooldown expired")
	}
}

func TestRemaining_ZeroWhenFresh(t *testing.T) {
	s := New(time.Minute)
	if r := s.Remaining("svc-a"); r != 0 {
		t.Fatalf("expected 0 remaining for unknown service, got %v", r)
	}
}

func TestRemaining_PositiveDuringCooldown(t *testing.T) {
	s := New(time.Minute)
	s.Allow("svc-a") //nolint
	if r := s.Remaining("svc-a"); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestAllow_EmptyServiceReturnsError(t *testing.T) {
	s := New(time.Minute)
	_, err := s.Allow("")
	if err == nil {
		t.Fatal("expected error for empty service name")
	}
}

func TestReset_EmptyServiceReturnsError(t *testing.T) {
	s := New(time.Minute)
	if err := s.Reset(""); err == nil {
		t.Fatal("expected error for empty service name")
	}
}
