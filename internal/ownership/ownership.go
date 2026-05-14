package ownership

import (
	"errors"
	"sync"
)

// Owner holds metadata about who owns a service.
type Owner struct {
	Service string `json:"service"`
	Team    string `json:"team"`
	Email   string `json:"email"`
	Slack   string `json:"slack,omitempty"`
}

// Store maps service names to their owners.
type Store struct {
	mu   sync.RWMutex
	data map[string]Owner
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[string]Owner)}
}

// Set registers or replaces the owner for a service.
func (s *Store) Set(o Owner) error {
	if o.Service == "" {
		return errors.New("service name is required")
	}
	if o.Team == "" {
		return errors.New("team is required")
	}
	if o.Email == "" {
		return errors.New("email is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[o.Service] = o
	return nil
}

// Get returns the owner for the given service.
func (s *Store) Get(service string) (Owner, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	o, ok := s.data[service]
	return o, ok
}

// Delete removes the owner record for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[service]; !ok {
		return errors.New("service not found")
	}
	delete(s.data, service)
	return nil
}

// All returns a snapshot of every owner record.
func (s *Store) All() []Owner {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Owner, 0, len(s.data))
	for _, o := range s.data {
		out = append(out, o)
	}
	return out
}
