// Package cooldown enforces a minimum quiet period between consecutive
// check triggers for a given service, preventing alert storms.
package cooldown

import (
	"errors"
	"sync"
	"time"
)

// ErrEmptyService is returned when an empty service name is provided.
var ErrEmptyService = errors.New("cooldown: service name must not be empty")

// entry tracks the last trigger time and cooldown duration for a service.
type entry struct {
	lastTriggered time.Time
	duration      time.Duration
}

// Store holds per-service cooldown state.
type Store struct {
	mu      sync.Mutex
	entries map[string]*entry
	default_ time.Duration
}

// New creates a Store with the given default cooldown duration.
func New(defaultDuration time.Duration) *Store {
	return &Store{
		entries:  make(map[string]*entry),
		default_: defaultDuration,
	}
}

// Set configures a custom cooldown duration for the named service.
func (s *Store) Set(service string, d time.Duration) error {
	if service == "" {
		return ErrEmptyService
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[service]; ok {
		e.duration = d
	} else {
		s.entries[service] = &entry{duration: d}
	}
	return nil
}

// Allow reports whether the service is outside its cooldown window.
// If allowed, the last-triggered timestamp is updated to now.
func (s *Store) Allow(service string) (bool, error) {
	if service == "" {
		return false, ErrEmptyService
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	e, ok := s.entries[service]
	if !ok {
		e = &entry{duration: s.default_}
		s.entries[service] = e
	}
	if !e.lastTriggered.IsZero() && now.Sub(e.lastTriggered) < e.duration {
		return false, nil
	}
	e.lastTriggered = now
	return true, nil
}

// Reset clears the last-triggered timestamp for a service, allowing an
// immediate re-trigger regardless of the cooldown window.
func (s *Store) Reset(service string) error {
	if service == "" {
		return ErrEmptyService
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[service]; ok {
		e.lastTriggered = time.Time{}
	}
	return nil
}

// Remaining returns the time left in the current cooldown window for the
// given service, or zero if the service is not in cooldown.
func (s *Store) Remaining(service string) time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[service]
	if !ok || e.lastTriggered.IsZero() {
		return 0
	}
	elapsed := time.Since(e.lastTriggered)
	if elapsed >= e.duration {
		return 0
	}
	return e.duration - elapsed
}
