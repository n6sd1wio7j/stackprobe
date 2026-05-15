package catalog

import (
	"errors"
	"sync"
	"time"
)

// Entry represents a registered service in the catalog.
type Entry struct {
	Service     string    `json:"service"`
	Description string    `json:"description"`
	Owner       string    `json:"owner"`
	URL         string    `json:"url"`
	Registered  time.Time `json:"registered"`
	Updated     time.Time `json:"updated"`
}

// Store holds catalog entries for all known services.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Entry
}

// New returns an empty catalog Store.
func New() *Store {
	return &Store{entries: make(map[string]Entry)}
}

// Register adds or updates a service entry in the catalog.
func (s *Store) Register(service, description, owner, url string) error {
	if service == "" {
		return errors.New("service name is required")
	}
	if url == "" {
		return errors.New("url is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	if existing, ok := s.entries[service]; ok {
		existing.Description = description
		existing.Owner = owner
		existing.URL = url
		existing.Updated = now
		s.entries[service] = existing
	} else {
		s.entries[service] = Entry{
			Service:     service,
			Description: description,
			Owner:       owner,
			URL:         url,
			Registered:  now,
			Updated:     now,
		}
	}
	return nil
}

// Get returns the catalog entry for a service.
func (s *Store) Get(service string) (Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// Delete removes a service from the catalog.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[service]; !ok {
		return errors.New("service not found")
	}
	delete(s.entries, service)
	return nil
}

// All returns a snapshot of all catalog entries.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
