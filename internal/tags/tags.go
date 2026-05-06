// Package tags provides service tagging and filtering capabilities.
package tags

import "sync"

// Store holds tag associations for services.
type Store struct {
	mu   sync.RWMutex
	tags map[string][]string // service name -> tags
}

// New creates a new tag Store.
func New() *Store {
	return &Store{
		tags: make(map[string][]string),
	}
}

// Set replaces all tags for a service.
func (s *Store) Set(service string, tags []string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copy := make([]string, len(tags))
	for i, t := range tags {
		copy[i] = t
	}
	s.tags[service] = copy
}

// Get returns the tags for a service. Returns nil if not found.
func (s *Store) Get(service string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tags[service]
	if !ok {
		return nil
	}
	out := make([]string, len(t))
	copy(out, t)
	return out
}

// Filter returns the names of services that have ALL of the given tags.
func (s *Store) Filter(required ...string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []string
serviceLoop:
	for svc, tags := range s.tags {
		for _, req := range required {
			if !contains(tags, req) {
				continue serviceLoop
			}
		}
		result = append(result, svc)
	}
	return result
}

// All returns a map of all service->tags associations.
func (s *Store) All() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]string, len(s.tags))
	for k, v := range s.tags {
		cp := make([]string, len(v))
		copy(cp, v)
		out[k] = cp
	}
	return out
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
