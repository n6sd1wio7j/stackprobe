// Package pinned tracks services that have been manually pinned to a
// specific status, overriding the result of live health checks.
package pinned

import (
	"errors"
	"sync"
	"time"
)

// Status represents a manually pinned health status.
type Status struct {
	Service  string    `json:"service"`
	Up       bool      `json:"up"`
	Reason   string    `json:"reason"`
	PinnedAt time.Time `json:"pinned_at"`
}

// Store holds pinned statuses for services.
type Store struct {
	mu    sync.RWMutex
	items map[string]Status
}

// New returns an initialised Store.
func New() *Store {
	return &Store{items: make(map[string]Status)}
}

// Pin records a manual status override for the given service.
func (s *Store) Pin(service string, up bool, reason string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[service] = Status{
		Service:  service,
		Up:       up,
		Reason:   reason,
		PinnedAt: time.Now().UTC(),
	}
	return nil
}

// Unpin removes any manual override for the given service.
func (s *Store) Unpin(service string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[service]; !ok {
		return errors.New("service not pinned")
	}
	delete(s.items, service)
	return nil
}

// Get returns the pinned status for a service and whether it exists.
func (s *Store) Get(service string) (Status, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.items[service]
	return v, ok
}

// IsPinned reports whether the service has a manual override.
func (s *Store) IsPinned(service string) bool {
	_, ok := s.Get(service)
	return ok
}

// All returns a snapshot of every pinned service.
func (s *Store) All() []Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Status, 0, len(s.items))
	for _, v := range s.items {
		out = append(out, v)
	}
	return out
}
