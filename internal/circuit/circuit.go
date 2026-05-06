package circuit

import (
	"sync"
	"time"
)

// State represents the circuit breaker state.
type State int

const (
	StateClosed State = iota
	StateOpen
	StateHalfOpen
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker is a per-service circuit breaker.
type Breaker struct {
	mu           sync.Mutex
	state        State
	failures      int
	threshold     int
	openUntil     time.Time
	cooldown      time.Duration
}

// Store holds a Breaker per service name.
type Store struct {
	mu       sync.Mutex
	breakers map[string]*Breaker
	threshold int
	cooldown  time.Duration
}

// New creates a new Store with the given failure threshold and cooldown window.
func New(threshold int, cooldown time.Duration) *Store {
	return &Store{
		breakers:  make(map[string]*Breaker),
		threshold: threshold,
		cooldown:  cooldown,
	}
}

func (s *Store) get(service string) *Breaker {
	s.mu.Lock()
	defer s.mu.Unlock()
	if b, ok := s.breakers[service]; ok {
		return b
	}
	b := &Breaker{threshold: s.threshold, cooldown: s.cooldown}
	s.breakers[service] = b
	return b
}

// Allow reports whether a request for the given service should proceed.
func (s *Store) Allow(service string) bool {
	b := s.get(service)
	b.mu.Lock()
	defer b.mu.Unlock()
	switch b.state {
	case StateOpen:
		if time.Now().After(b.openUntil) {
			b.state = StateHalfOpen
			return true
		}
		return false
	default:
		return true
	}
}

// RecordSuccess records a successful check, closing the circuit if half-open.
func (s *Store) RecordSuccess(service string) {
	b := s.get(service)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed check, opening the circuit after threshold.
func (s *Store) RecordFailure(service string) {
	b := s.get(service)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.failures++
	if b.failures >= b.threshold {
		b.state = StateOpen
		b.openUntil = time.Now().Add(b.cooldown)
	}
}

// State returns the current state of the breaker for the given service.
func (s *Store) State(service string) State {
	b := s.get(service)
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}
