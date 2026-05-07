// Package maintenance provides a way to mark services as under maintenance,
// suppressing alerts and health-check failures during planned downtime.
package maintenance

import (
	"sync"
	"time"
)

// Window describes a maintenance period for a single service.
type Window struct {
	Service   string    `json:"service"`
	Reason    string    `json:"reason"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
}

// Active reports whether the window covers the given instant.
func (w Window) Active(at time.Time) bool {
	return !at.Before(w.StartsAt) && at.Before(w.EndsAt)
}

// Store holds maintenance windows keyed by service name.
type Store struct {
	mu      sync.RWMutex
	windows map[string]Window
}

// New returns an initialised Store.
func New() *Store {
	return &Store{windows: make(map[string]Window)}
}

// Set registers or replaces the maintenance window for a service.
func (s *Store) Set(w Window) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[w.Service] = w
}

// Delete removes the maintenance window for a service.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.windows, service)
}

// IsActive returns true when the named service has an active maintenance
// window at the given instant.
func (s *Store) IsActive(service string, at time.Time) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.windows[service]
	if !ok {
		return false
	}
	return w.Active(at)
}

// All returns a snapshot of every registered window.
func (s *Store) All() []Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Window, 0, len(s.windows))
	for _, w := range s.windows {
		out = append(out, w)
	}
	return out
}
