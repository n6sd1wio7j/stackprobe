package passthrough

import (
	"errors"
	"sync"
)

// Mode represents whether a service is in passthrough mode.
type Mode int

const (
	ModeNormal    Mode = iota // health checks run normally
	ModePassthrough           // always report healthy, skip checks
)

// Store holds per-service passthrough overrides.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Mode
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Mode)}
}

// Set configures the passthrough mode for a service.
func (s *Store) Set(service string, m Mode) error {
	if service == "" {
		return errors.New("passthrough: service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = m
	return nil
}

// Get returns the Mode for a service. Unknown services return ModeNormal.
func (s *Store) Get(service string) Mode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if m, ok := s.entries[service]; ok {
		return m
	}
	return ModeNormal
}

// IsPassthrough reports whether the named service is currently in passthrough mode.
func (s *Store) IsPassthrough(service string) bool {
	return s.Get(service) == ModePassthrough
}

// Delete removes any override for a service, reverting it to ModeNormal.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, service)
}

// All returns a snapshot of every service and its current mode.
func (s *Store) All() map[string]Mode {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Mode, len(s.entries))
	for k, v := range s.entries {
		out[k] = v
	}
	return out
}
