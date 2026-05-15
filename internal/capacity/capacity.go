package capacity

import (
	"errors"
	"sync"
)

// Record holds capacity information for a service.
type Record struct {
	Service  string  `json:"service"`
	Limit    int     `json:"limit"`
	Current  int     `json:"current"`
	Percent  float64 `json:"percent"`
}

// Store tracks capacity usage per service.
type Store struct {
	mu      sync.RWMutex
	entries map[string]Record
}

// New returns an initialised Store.
func New() *Store {
	return &Store{entries: make(map[string]Record)}
}

// Set registers or updates capacity for a service.
func (s *Store) Set(service string, limit, current int) error {
	if service == "" {
		return errors.New("service name required")
	}
	if limit <= 0 {
		return errors.New("limit must be greater than zero")
	}
	if current < 0 {
		return errors.New("current must be non-negative")
	}
	pct := float64(current) / float64(limit) * 100
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = Record{
		Service: service,
		Limit:   limit,
		Current: current,
		Percent: pct,
	}
	return nil
}

// Get returns the capacity record for a service.
func (s *Store) Get(service string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.entries[service]
	return r, ok
}

// Delete removes the capacity record for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[service]; !ok {
		return errors.New("unknown service")
	}
	delete(s.entries, service)
	return nil
}

// All returns a snapshot of all capacity records.
func (s *Store) All() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, 0, len(s.entries))
	for _, r := range s.entries {
		out = append(out, r)
	}
	return out
}

// IsSaturated reports whether current usage meets or exceeds the limit.
func (s *Store) IsSaturated(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.entries[service]
	if !ok {
		return false
	}
	return r.Current >= r.Limit
}
