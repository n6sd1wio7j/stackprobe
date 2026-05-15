package checkpoint

import (
	"errors"
	"sync"
	"time"
)

// Record holds the last known good state for a service.
type Record struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	SavedAt   time.Time `json:"saved_at"`
	Note      string    `json:"note,omitempty"`
}

// Store persists checkpoint records keyed by service name.
type Store struct {
	mu      sync.RWMutex
	records map[string]Record
}

// New returns an empty Store.
func New() *Store {
	return &Store{records: make(map[string]Record)}
}

// Save creates or overwrites the checkpoint for the given service.
func (s *Store) Save(service, status, note string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if status == "" {
		return errors.New("status must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[service] = Record{
		Service: service,
		Status:  status,
		Note:    note,
		SavedAt: time.Now().UTC(),
	}
	return nil
}

// Get returns the checkpoint for a service, or an error if not found.
func (s *Store) Get(service string) (Record, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.records[service]
	if !ok {
		return Record{}, errors.New("no checkpoint for service: " + service)
	}
	return r, nil
}

// Delete removes the checkpoint for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[service]; !ok {
		return errors.New("no checkpoint for service: " + service)
	}
	delete(s.records, service)
	return nil
}

// All returns a snapshot of every stored checkpoint.
func (s *Store) All() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, 0, len(s.records))
	for _, r := range s.records {
		out = append(out, r)
	}
	return out
}
