// Package dependency tracks inter-service dependencies and resolves
// transitive health based on upstream status.
package dependency

import (
	"errors"
	"sync"
)

// ErrCycle is returned when adding a dependency would create a cycle.
var ErrCycle = errors.New("dependency: cycle detected")

// Store holds a directed dependency graph between services.
type Store struct {
	mu    sync.RWMutex
	edges map[string][]string // service -> list of upstream dependencies
}

// New creates an empty dependency Store.
func New() *Store {
	return &Store{edges: make(map[string][]string)}
}

// Add registers that service depends on upstream.
// Returns ErrCycle if the relationship would introduce a cycle.
func (s *Store) Add(service, upstream string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check that upstream does not already (transitively) depend on service.
	if s.reachable(upstream, service) {
		return ErrCycle
	}

	s.edges[service] = append(s.edges[service], upstream)
	return nil
}

// Deps returns the direct dependencies of service.
func (s *Store) Deps(service string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]string, len(s.edges[service]))
	copy(out, s.edges[service])
	return out
}

// Affected returns all services that (directly or transitively) depend on
// upstream — i.e. services that may be impacted when upstream is unhealthy.
func (s *Store) Affected(upstream string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []string
	visited := make(map[string]bool)
	s.collectAffected(upstream, visited, &result)
	return result
}

// Remove deletes all dependency edges for service.
func (s *Store) Remove(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.edges, service)
}

// reachable reports whether target is reachable from start (caller must hold lock).
func (s *Store) reachable(start, target string) bool {
	if start == target {
		return true
	}
	for _, dep := range s.edges[start] {
		if s.reachable(dep, target) {
			return true
		}
	}
	return false
}

// collectAffected performs a reverse DFS to find dependents of upstream.
func (s *Store) collectAffected(upstream string, visited map[string]bool, result *[]string) {
	for svc, deps := range s.edges {
		if visited[svc] {
			continue
		}
		for _, d := range deps {
			if d == upstream {
				visited[svc] = true
				*result = append(*result, svc)
				s.collectAffected(svc, visited, result)
				break
			}
		}
	}
}
