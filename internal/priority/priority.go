package priority

import (
	"errors"
	"sync"
)

// Level represents a service priority level.
type Level int

const (
	Low    Level = 1
	Medium Level = 2
	High   Level = 3
	Critical Level = 4
)

var levelNames = map[Level]string{
	Low:      "low",
	Medium:   "medium",
	High:     "high",
	Critical: "critical",
}

func (l Level) String() string {
	if name, ok := levelNames[l]; ok {
		return name
	}
	return "unknown"
}

// Store holds priority levels keyed by service name.
type Store struct {
	mu       sync.RWMutex
	levels   map[string]Level
	defLevel Level
}

// New creates a new Store with the given default level.
func New(defaultLevel Level) *Store {
	return &Store{
		levels:   make(map[string]Level),
		defLevel: defaultLevel,
	}
}

// Set assigns a priority level to a service.
func (s *Store) Set(service string, level Level) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if _, ok := levelNames[level]; !ok {
		return errors.New("invalid priority level")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.levels[service] = level
	return nil
}

// Get returns the priority level for a service, falling back to the default.
func (s *Store) Get(service string) Level {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if l, ok := s.levels[service]; ok {
		return l
	}
	return s.defLevel
}

// Delete removes any explicit priority override for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.levels[service]; !ok {
		return errors.New("service not found")
	}
	delete(s.levels, service)
	return nil
}

// All returns a snapshot of all explicitly set priorities.
func (s *Store) All() map[string]Level {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]Level, len(s.levels))
	for k, v := range s.levels {
		out[k] = v
	}
	return out
}

// Filter returns service names whose priority is at or above minLevel.
func (s *Store) Filter(minLevel Level) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []string
	for svc, lvl := range s.levels {
		if lvl >= minLevel {
			result = append(result, svc)
		}
	}
	return result
}
