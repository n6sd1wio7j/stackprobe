package grouping

import (
	"errors"
	"sync"
)

// Group represents a named collection of service names.
type Group struct {
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

// Store manages service groupings.
type Store struct {
	mu     sync.RWMutex
	groups map[string][]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{groups: make(map[string][]string)}
}

// Set creates or replaces the group identified by name.
func (s *Store) Set(name string, services []string) error {
	if name == "" {
		return errors.New("group name must not be empty")
	}
	if len(services) == 0 {
		return errors.New("services must not be empty")
	}
	copy := make([]string, len(services))
	for i, svc := range services {
		if svc == "" {
			return errors.New("service name must not be empty")
		}
		copy[i] = svc
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.groups[name] = copy
	return nil
}

// Get returns the services belonging to the named group.
func (s *Store) Get(name string) ([]string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	svcs, ok := s.groups[name]
	if !ok {
		return nil, false
	}
	out := make([]string, len(svcs))
	copy(out, svcs)
	return out, true
}

// Delete removes a group by name.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.groups[name]; !ok {
		return errors.New("group not found")
	}
	delete(s.groups, name)
	return nil
}

// All returns a snapshot of every group.
func (s *Store) All() []Group {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Group, 0, len(s.groups))
	for name, svcs := range s.groups {
		copy := make([]string, len(svcs))
		for i, v := range svcs {
			copy[i] = v
		}
		out = append(out, Group{Name: name, Services: copy})
	}
	return out
}

// Member reports whether service belongs to the named group.
func (s *Store) Member(name, service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, svc := range s.groups[name] {
		if svc == service {
			return true
		}
	}
	return false
}
