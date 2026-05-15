package lockout

import (
	"errors"
	"sync"
	"time"
)

// ErrAlreadyLocked is returned when a service is already locked out.
var ErrAlreadyLocked = errors.New("service is already locked out")

// ErrNotLocked is returned when attempting to unlock a service that is not locked.
var ErrNotLocked = errors.New("service is not locked out")

// ErrEmptyService is returned when an empty service name is provided.
var ErrEmptyService = errors.New("service name must not be empty")

// Entry represents a lockout record for a service.
type Entry struct {
	Service   string    `json:"service"`
	Reason    string    `json:"reason"`
	LockedAt  time.Time `json:"locked_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired reports whether the lockout window has passed.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Store holds lockout entries keyed by service name.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Entry)}
}

// Lock creates a lockout entry for the given service.
func (s *Store) Lock(service, reason string, duration time.Duration) error {
	if service == "" {
		return ErrEmptyService
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok := s.entries[service]; ok && !e.IsExpired() {
		return ErrAlreadyLocked
	}
	now := time.Now()
	s.entries[service] = Entry{
		Service:   service,
		Reason:    reason,
		LockedAt:  now,
		ExpiresAt: now.Add(duration),
	}
	return nil
}

// Unlock removes the lockout entry for the given service.
func (s *Store) Unlock(service string) error {
	if service == "" {
		return ErrEmptyService
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[service]; !ok {
		return ErrNotLocked
	}
	delete(s.entries, service)
	return nil
}

// IsLocked reports whether the service is currently locked out.
func (s *Store) IsLocked(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return ok && !e.IsExpired()
}

// Get returns the lockout entry for a service.
func (s *Store) Get(service string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// All returns a snapshot of all current (non-expired) lockout entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		if !e.IsExpired() {
			out = append(out, e)
		}
	}
	return out
}
