package timeout

import (
	"errors"
	"sync"
	"time"
)

// ErrTimeout is returned when a service-specific timeout override is not found.
var ErrTimeout = errors.New("timeout: no override set for service")

// Store holds per-service timeout overrides.
type Store struct {
	mu       sync.RWMutex
	defaults time.Duration
	overrides map[string]time.Duration
}

// New creates a Store with the given default timeout.
func New(defaultTimeout time.Duration) *Store {
	return &Store{
		defaults:  defaultTimeout,
		overrides: make(map[string]time.Duration),
	}
}

// Set stores a timeout override for the named service.
func (s *Store) Set(service string, d time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.overrides[service] = d
}

// Get returns the timeout for the named service.
// If no override exists the default is returned.
func (s *Store) Get(service string) time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if d, ok := s.overrides[service]; ok {
		return d
	}
	return s.defaults
}

// Delete removes a per-service override, reverting to the default.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.overrides, service)
}

// All returns a snapshot of every current override.
func (s *Store) All() map[string]time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]time.Duration, len(s.overrides))
	for k, v := range s.overrides {
		out[k] = v
	}
	return out
}
