// Package environ tracks per-service environment metadata (e.g. "production", "staging").
package environ

import (
	"errors"
	"sync"
)

// Store holds environment assignments for services.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Set assigns an environment label to a service.
func (s *Store) Set(service, env string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if env == "" {
		return errors.New("environment must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[service] = env
	return nil
}

// Get returns the environment for a service and whether it was found.
func (s *Store) Get(service string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	env, ok := s.data[service]
	return env, ok
}

// Delete removes the environment assignment for a service.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, service)
}

// Filter returns all services whose environment matches the given value.
func (s *Store) Filter(env string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []string
	for svc, e := range s.data {
		if e == env {
			out = append(out, svc)
		}
	}
	return out
}

// All returns a snapshot of every service→environment mapping.
func (s *Store) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out
}
