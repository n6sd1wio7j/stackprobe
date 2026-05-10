// Package flipper provides a simple feature-flag store for stackprobe.
// Flags can be toggled at runtime and queried by name.
package flipper

import (
	"fmt"
	"sync"
	"time"
)

// Flag represents a single feature flag.
type Flag struct {
	Name      string    `json:"name"`
	Enabled   bool      `json:"enabled"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store holds feature flags in memory.
type Store struct {
	mu    sync.RWMutex
	flags map[string]*Flag
}

// New returns an initialised Store.
func New() *Store {
	return &Store{flags: make(map[string]*Flag)}
}

// Set creates or updates a flag with the given enabled state.
func (s *Store) Set(name string, enabled bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flags[name] = &Flag{
		Name:      name,
		Enabled:   enabled,
		UpdatedAt: time.Now().UTC(),
	}
}

// Get returns the Flag for name and true, or an error if not found.
func (s *Store) Get(name string) (Flag, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flags[name]
	if !ok {
		return Flag{}, fmt.Errorf("flipper: flag %q not found", name)
	}
	return *f, nil
}

// IsEnabled returns true when the named flag exists and is enabled.
func (s *Store) IsEnabled(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	f, ok := s.flags[name]
	return ok && f.Enabled
}

// Delete removes a flag from the store.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.flags[name]; !ok {
		return fmt.Errorf("flipper: flag %q not found", name)
	}
	delete(s.flags, name)
	return nil
}

// All returns a snapshot of every flag in the store.
func (s *Store) All() []Flag {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Flag, 0, len(s.flags))
	for _, f := range s.flags {
		out = append(out, *f)
	}
	return out
}
