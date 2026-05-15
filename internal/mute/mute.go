package mute

import (
	"errors"
	"sync"
	"time"
)

// Window describes a mute period for a service.
type Window struct {
	Service string    `json:"service"`
	Reason  string    `json:"reason"`
	Until   time.Time `json:"until"`
}

// Store holds mute windows keyed by service name.
type Store struct {
	mu      sync.RWMutex
	windows map[string]Window
}

// New returns an initialised Store.
func New() *Store {
	return &Store{windows: make(map[string]Window)}
}

// Mute silences alerts for the given service until the specified time.
func (s *Store) Mute(service, reason string, until time.Time) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if until.IsZero() || !until.After(time.Now()) {
		return errors.New("until must be a future time")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.windows[service] = Window{Service: service, Reason: reason, Until: until}
	return nil
}

// Unmute removes an active mute for the given service.
func (s *Store) Unmute(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.windows[service]; !ok {
		return errors.New("no mute window found for service")
	}
	delete(s.windows, service)
	return nil
}

// IsMuted reports whether the service is currently muted.
func (s *Store) IsMuted(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.windows[service]
	return ok && time.Now().Before(w.Until)
}

// Get returns the mute window for a service, if one exists.
func (s *Store) Get(service string) (Window, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	w, ok := s.windows[service]
	return w, ok
}

// All returns a snapshot of all currently active mute windows.
func (s *Store) All() []Window {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	out := make([]Window, 0, len(s.windows))
	for _, w := range s.windows {
		if now.Before(w.Until) {
			out = append(out, w)
		}
	}
	return out
}
