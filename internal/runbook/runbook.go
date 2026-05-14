package runbook

import (
	"errors"
	"sync"
	"time"
)

// Entry holds a runbook entry linked to a service.
type Entry struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

// Store manages runbook entries keyed by service name.
type Store struct {
	mu      sync.RWMutex
	entries map[string]*Entry
	nextID  int
}

// New returns an empty Store.
func New() *Store {
	return &Store{entries: make(map[string]*Entry)}
}

// Set creates or replaces the runbook entry for a service.
func (s *Store) Set(service, title, url, notes string) (*Entry, error) {
	if service == "" {
		return nil, errors.New("service must not be empty")
	}
	if title == "" {
		return nil, errors.New("title must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	e := &Entry{
		ID:        fmt.Sprintf("%d", s.nextID),
		Service:   service,
		Title:     title,
		URL:       url,
		Notes:     notes,
		CreatedAt: time.Now().UTC(),
	}
	s.entries[service] = e
	return e, nil
}

// Get returns the runbook entry for a service.
func (s *Store) Get(service string) (*Entry, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[service]
	return e, ok
}

// Delete removes the runbook entry for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[service]; !ok {
		return errors.New("runbook entry not found")
	}
	delete(s.entries, service)
	return nil
}

// All returns a snapshot of all entries.
func (s *Store) All() []*Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, e)
	}
	return out
}
