package deadletter

import (
	"errors"
	"sync"
	"time"
)

// ErrQueueFull is returned when the dead-letter queue has reached capacity.
var ErrQueueFull = errors.New("dead-letter queue is full")

// ErrUnknownService is returned when no entries exist for a service.
var ErrUnknownService = errors.New("unknown service")

// Entry represents a single failed check that could not be processed.
type Entry struct {
	Service   string    `json:"service"`
	Reason    string    `json:"reason"`
	OccurredAt time.Time `json:"occurred_at"`
}

// Store holds dead-letter entries per service with a configurable capacity.
type Store struct {
	mu       sync.RWMutex
	entries  map[string][]Entry
	capacity int
}

// New creates a Store with the given per-service capacity.
func New(capacity int) *Store {
	if capacity <= 0 {
		capacity = 100
	}
	return &Store{
		entries:  make(map[string][]Entry),
		capacity: capacity,
	}
}

// Push appends a failed-check entry for the given service.
// Returns ErrQueueFull when the per-service limit is reached.
func (s *Store) Push(service, reason string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.entries[service]) >= s.capacity {
		return ErrQueueFull
	}
	s.entries[service] = append(s.entries[service], Entry{
		Service:    service,
		Reason:     reason,
		OccurredAt: time.Now().UTC(),
	})
	return nil
}

// Get returns all dead-letter entries for a service.
func (s *Store) Get(service string) ([]Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	ent, ok := s.entries[service]
	if !ok {
		return nil, ErrUnknownService
	}
	out := make([]Entry, len(ent))
	copy(out, ent)
	return out, nil
}

// Flush removes all entries for a service and returns them.
func (s *Store) Flush(service string) ([]Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	ent, ok := s.entries[service]
	if !ok {
		return nil, ErrUnknownService
	}
	delete(s.entries, service)
	return ent, nil
}

// All returns a snapshot of every entry across all services.
func (s *Store) All() map[string][]Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]Entry, len(s.entries))
	for svc, ents := range s.entries {
		cp := make([]Entry, len(ents))
		copy(cp, ents)
		out[svc] = cp
	}
	return out
}
