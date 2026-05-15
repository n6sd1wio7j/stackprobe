package endpoint

import (
	"errors"
	"sync"
	"time"
)

// Meta holds metadata about a registered endpoint.
type Meta struct {
	Service     string    `json:"service"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	RegisteredAt time.Time `json:"registered_at"`
}

// Store manages endpoint registrations.
type Store struct {
	mu    sync.RWMutex
	items map[string]Meta
}

// New returns an initialised Store.
func New() *Store {
	return &Store{items: make(map[string]Meta)}
}

// Register adds or replaces an endpoint entry.
func (s *Store) Register(service, url, description string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if url == "" {
		return errors.New("url must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[service] = Meta{
		Service:      service,
		URL:          url,
		Description:  description,
		RegisteredAt: time.Now().UTC(),
	}
	return nil
}

// Get returns the Meta for a service, or an error if not found.
func (s *Store) Get(service string) (Meta, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.items[service]
	if !ok {
		return Meta{}, errors.New("endpoint not found: " + service)
	}
	return m, nil
}

// Delete removes the endpoint for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[service]; !ok {
		return errors.New("endpoint not found: " + service)
	}
	delete(s.items, service)
	return nil
}

// All returns a snapshot of all registered endpoints.
func (s *Store) All() []Meta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Meta, 0, len(s.items))
	for _, m := range s.items {
		out = append(out, m)
	}
	return out
}
