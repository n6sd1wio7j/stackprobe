package depgraph

import (
	"errors"
	"sync"
)

// Node represents a service node in the dependency graph with a health status.
type Node struct {
	Service  string
	Healthy  bool
	Children []string
}

// Store holds a directed graph of service dependencies and their health.
type Store struct {
	mu    sync.RWMutex
	nodes map[string]*Node
}

// New returns an initialised Store.
func New() *Store {
	return &Store{nodes: make(map[string]*Node)}
}

// Register adds or replaces a service node with its direct dependencies.
func (s *Store) Register(service string, deps []string) error {
	if service == "" {
		return errors.New("service name required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes[service] = &Node{Service: service, Healthy: true, Children: deps}
	return nil
}

// SetHealth updates the healthy flag for a service.
func (s *Store) SetHealth(service string, healthy bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	n, ok := s.nodes[service]
	if !ok {
		return errors.New("unknown service: " + service)
	}
	n.Healthy = healthy
	return nil
}

// Get returns the node for a service.
func (s *Store) Get(service string) (*Node, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n, ok := s.nodes[service]
	if !ok {
		return nil, false
	}
	copy := *n
	return &copy, true
}

// All returns all registered nodes.
func (s *Store) All() []*Node {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Node, 0, len(s.nodes))
	for _, n := range s.nodes {
		copy := *n
		out = append(out, &copy)
	}
	return out
}

// Unhealthy returns the names of services that are currently unhealthy.
func (s *Store) Unhealthy() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []string
	for _, n := range s.nodes {
		if !n.Healthy {
			out = append(out, n.Service)
		}
	}
	return out
}
