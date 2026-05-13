package weight

import (
	"errors"
	"sync"
)

// Store holds load-balancing weights for services.
type Store struct {
	mu      sync.RWMutex
	weights map[string]int
	default_ int
}

// New returns a Store with the given default weight.
func New(defaultWeight int) *Store {
	if defaultWeight <= 0 {
		defaultWeight = 1
	}
	return &Store{
		weights:  make(map[string]int),
		default_: defaultWeight,
	}
}

// Set assigns a weight to a service. Weight must be >= 1.
func (s *Store) Set(service string, weight int) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if weight < 1 {
		return errors.New("weight must be at least 1")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weights[service] = weight
	return nil
}

// Get returns the weight for a service, falling back to the default.
func (s *Store) Get(service string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if w, ok := s.weights[service]; ok {
		return w
	}
	return s.default_
}

// Delete removes an explicit weight, reverting to the default.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.weights, service)
}

// All returns a snapshot of all explicitly set weights.
func (s *Store) All() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]int, len(s.weights))
	for k, v := range s.weights {
		out[k] = v
	}
	return out
}

// Default returns the fallback weight.
func (s *Store) Default() int {
	return s.default_
}
