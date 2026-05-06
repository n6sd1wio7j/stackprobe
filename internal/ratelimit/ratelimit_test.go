package ratelimit_test

import (
	"testing"
	"time"

	"github.com/example/stackprobe/internal/ratelimit"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	if !l.Allow("svc-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallBlocked(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	l.Allow("svc-a")
	if l.Allow("svc-a") {
		t.Fatal("expected second immediate call to be blocked")
	}
}

func TestAllow_AfterIntervalPermitted(t *testing.T) {
	l := ratelimit.New(50 * time.Millisecond)
	l.Allow("svc-b")
	time.Sleep(60 * time.Millisecond)
	if !l.Allow("svc-b") {
		t.Fatal("expected call after interval to be allowed")
	}
}

func TestAllow_IndependentServices(t *testing.T) {
	l := ratelimit.New(500 * time.Millisecond)
	l.Allow("svc-x")
	if !l.Allow("svc-y") {
		t.Fatal("expected different service to be allowed independently")
	}
}

func TestReset_AllowsImmediateRecheck(t *testing.T) {
	l := ratelimit.New(500 * time.Millisecond)
	l.Allow("svc-c")
	l.Reset("svc-c")
	if !l.Allow("svc-c") {
		t.Fatal("expected allow after reset")
	}
}

func TestNextAllowed_UnknownServiceZeroTime(t *testing.T) {
	l := ratelimit.New(100 * time.Millisecond)
	if !l.NextAllowed("unknown").IsZero() {
		t.Fatal("expected zero time for unknown service")
	}
}

func TestNextAllowed_KnownService(t *testing.T) {
	interval := 200 * time.Millisecond
	l := ratelimit.New(interval)
	before := time.Now()
	l.Allow("svc-d")
	next := l.NextAllowed("svc-d")
	if next.Before(before.Add(interval)) {
		t.Fatalf("expected next allowed >= %v, got %v", before.Add(interval), next)
	}
}

func TestNew_ZeroIntervalDefaultsToOneSecond(t *testing.T) {
	l := ratelimit.New(0)
	l.Allow("svc-e")
	next := l.NextAllowed("svc-e")
	if next.IsZero() {
		t.Fatal("expected non-zero next allowed time")
	}
	if time.Until(next) < 900*time.Millisecond {
		t.Fatal("expected interval to default to ~1 second")
	}
}
