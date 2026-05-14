package endpointmeta

import (
	"errors"
	"sync"
)

// Meta holds descriptive metadata for a single endpoint.
type Meta struct {
	Service     string `json:"service"`
	Description string `json:"description"`
	Owner       string `json:"owner"`
	DocURL      string `json:"doc_url"`
	Tier        string `json:"tier"`
}

// Store holds endpoint metadata keyed by service name.
type Store struct {
	mu    sync.RWMutex
	items map[string]Meta
}

// New returns an initialised Store.
func New() *Store {
	return &Store{items: make(map[string]Meta)}
}

// Set stores or replaces the metadata for the given service.
func (s *Store) Set(service string, m Meta) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	m.Service = service
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[service] = m
	return nil
}

// Get returns the metadata for the given service.
func (s *Store) Get(service string) (Meta, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.items[service]
	return m, ok
}

// Delete removes the metadata entry for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.items[service]; !ok {
		return errors.New("service not found")
	}
	delete(s.items, service)
	return nil
}

// All returns a snapshot of every stored metadata entry.
func (s *Store) All() []Meta {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Meta, 0, len(s.items))
	for _, m := range s.items {
		out = append(out, m)
	}
	return out
}
