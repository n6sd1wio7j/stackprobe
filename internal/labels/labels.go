package labels

import (
	"fmt"
	"sync"
)

// Store holds arbitrary key-value label pairs for each service.
type Store struct {
	mu   sync.RWMutex
	data map[string]map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[string]map[string]string)}
}

// Set assigns a label key/value for the given service.
func (s *Store) Set(service, key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.data[service]; !ok {
		s.data[service] = make(map[string]string)
	}
	s.data[service][key] = value
}

// Get returns all labels for a service, or an error if unknown.
func (s *Store) Get(service string) (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	labels, ok := s.data[service]
	if !ok {
		return nil, fmt.Errorf("labels: unknown service %q", service)
	}
	copy := make(map[string]string, len(labels))
	for k, v := range labels {
		copy[k] = v
	}
	return copy, nil
}

// Delete removes a single label key from a service.
func (s *Store) Delete(service, key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if m, ok := s.data[service]; ok {
		delete(m, key)
	}
}

// Filter returns all services whose labels contain all of the supplied
// key-value pairs.
func (s *Store) Filter(selector map[string]string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var matched []string
outer:
	for svc, labels := range s.data {
		for k, v := range selector {
			if labels[k] != v {
				continue outer
			}
		}
		matched = append(matched, svc)
	}
	return matched
}

// All returns a snapshot of every service and its labels.
func (s *Store) All() map[string]map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]map[string]string, len(s.data))
	for svc, labels := range s.data {
		copy := make(map[string]string, len(labels))
		for k, v := range labels {
			copy[k] = v
		}
		out[svc] = copy
	}
	return out
}
