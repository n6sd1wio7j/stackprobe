package envelope

import (
	"errors"
	"sync"
	"time"
)

// Envelope wraps a health-check result with metadata for transport.
type Envelope struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Message   string    `json:"message,omitempty"`
	WrappedAt time.Time `json:"wrapped_at"`
	TTL       int       `json:"ttl_seconds"`
}

// Store holds envelopes keyed by service name.
type Store struct {
	mu      sync.RWMutex
	items   map[string]Envelope
	defaultTTL int
}

// New returns a Store with the given default TTL in seconds.
func New(defaultTTL int) *Store {
	if defaultTTL <= 0 {
		defaultTTL = 60
	}
	return &Store{
		items:      make(map[string]Envelope),
		defaultTTL: defaultTTL,
	}
}

// Wrap creates or replaces the envelope for a service.
func (s *Store) Wrap(service, status, message string) error {
	if service == "" {
		return errors.New("service must not be empty")
	}
	if status == "" {
		return errors.New("status must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[service] = Envelope{
		Service:   service,
		Status:    status,
		Message:   message,
		WrappedAt: time.Now().UTC(),
		TTL:       s.defaultTTL,
	}
	return nil
}

// Get returns the envelope for a service.
func (s *Store) Get(service string) (Envelope, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.items[service]
	return e, ok
}

// Delete removes the envelope for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[service]; !ok {
		return errors.New("unknown service")
	}
	delete(s.items, service)
	return nil
}

// All returns a snapshot of all envelopes.
func (s *Store) All() []Envelope {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Envelope, 0, len(s.items))
	for _, e := range s.items {
		out = append(out, e)
	}
	return out
}

// IsExpired reports whether the envelope has exceeded its TTL.
func (s *Store) IsExpired(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.items[service]
	if !ok {
		return true
	}
	return time.Since(e.WrappedAt) > time.Duration(e.TTL)*time.Second
}
