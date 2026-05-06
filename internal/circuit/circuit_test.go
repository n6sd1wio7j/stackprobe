package circuit_test

import (
	"testing"
	"time"

	"github.com/user/stackprobe/internal/circuit"
)

const (
	svc       = "api"
	threshold = 3
	cooldown  = 100 * time.Millisecond
)

func newStore() *circuit.Store {
	return circuit.New(threshold, cooldown)
}

// openCircuit is a helper that trips the circuit breaker for the given service
// by recording the configured number of failures.
func openCircuit(s *circuit.Store, service string) {
	for i := 0; i < threshold; i++ {
		s.RecordFailure(service)
	}
}

func TestAllow_InitiallyPermits(t *testing.T) {
	s := newStore()
	if !s.Allow(svc) {
		t.Fatal("expected Allow=true for fresh service")
	}
}

func TestState_InitiallyClosed(t *testing.T) {
	s := newStore()
	if got := s.State(svc); got != circuit.StateClosed {
		t.Fatalf("expected closed, got %s", got)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	s := newStore()
	openCircuit(s, svc)
	if got := s.State(svc); got != circuit.StateOpen {
		t.Fatalf("expected open after %d failures, got %s", threshold, got)
	}
	if s.Allow(svc) {
		t.Fatal("expected Allow=false when circuit is open")
	}
}

func TestRecordSuccess_ResetsClosed(t *testing.T) {
	s := newStore()
	for i := 0; i < threshold-1; i++ {
		s.RecordFailure(svc)
	}
	s.RecordSuccess(svc)
	if got := s.State(svc); got != circuit.StateClosed {
		t.Fatalf("expected closed after success, got %s", got)
	}
	if !s.Allow(svc) {
		t.Fatal("expected Allow=true after reset")
	}
}

func TestHalfOpen_AfterCooldown(t *testing.T) {
	s := newStore()
	openCircuit(s, svc)
	time.Sleep(cooldown + 10*time.Millisecond)
	if !s.Allow(svc) {
		t.Fatal("expected Allow=true in half-open state after cooldown")
	}
	if got := s.State(svc); got != circuit.StateHalfOpen {
		t.Fatalf("expected half-open, got %s", got)
	}
}

func TestHalfOpen_SuccessCloses(t *testing.T) {
	s := newStore()
	openCircuit(s, svc)
	time.Sleep(cooldown + 10*time.Millisecond)
	s.Allow(svc) // transition to half-open
	s.RecordSuccess(svc)
	if got := s.State(svc); got != circuit.StateClosed {
		t.Fatalf("expected closed after half-open success, got %s", got)
	}
}

func TestIndependentServices(t *testing.T) {
	s := newStore()
	openCircuit(s, "broken")
	if !s.Allow("healthy") {
		t.Fatal("unrelated service should not be affected by another's failures")
	}
}
