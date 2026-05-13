// Package drain implements graceful drain management for services.
// A service in drain mode is excluded from health-check routing while
// in-flight requests are allowed to complete.
package drain

import (
	"errors"
	"sync"
	"time"
)

// ErrUnknownService is returned when no drain entry exists for a service.
var ErrUnknownService = errors.New("drain: unknown service")

// ErrAlreadyDraining is returned when the service is already draining.
var ErrAlreadyDraining = errors.New("drain: service already draining")

// State holds drain information for a single service.
type State struct {
	Service   string    `json:"service"`
	Draining  bool      `json:"draining"`
	StartedAt time.Time `json:"started_at,omitempty"`
}

// Store tracks drain state for services.
type Store struct {
	mu     sync.RWMutex
	states map[string]*State
}

// New returns an initialised Store.
func New() *Store {
	return &Store{states: make(map[string]*State)}
}

// Enable marks a service as draining. Returns ErrAlreadyDraining if it is
// already in drain mode.
func (s *Store) Enable(service string) error {
	if service == "" {
		return errors.New("drain: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if st, ok := s.states[service]; ok && st.Draining {
		return ErrAlreadyDraining
	}
	s.states[service] = &State{
		Service:   service,
		Draining:  true,
		StartedAt: time.Now().UTC(),
	}
	return nil
}

// Disable clears the drain flag for a service.
func (s *Store) Disable(service string) error {
	if service == "" {
		return errors.New("drain: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	st, ok := s.states[service]
	if !ok || !st.Draining {
		return ErrUnknownService
	}
	delete(s.states, service)
	return nil
}

// IsDraining reports whether a service is currently in drain mode.
func (s *Store) IsDraining(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.states[service]
	return ok && st.Draining
}

// Get returns the State for a service or ErrUnknownService.
func (s *Store) Get(service string) (State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.states[service]
	if !ok {
		return State{}, ErrUnknownService
	}
	return *st, nil
}

// All returns a snapshot of every service currently in drain mode.
func (s *Store) All() []State {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]State, 0, len(s.states))
	for _, st := range s.states {
		out = append(out, *st)
	}
	return out
}
