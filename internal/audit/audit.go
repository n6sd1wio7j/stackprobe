package audit

import (
	"fmt"
	"sync"
	"time"
)

// EventKind classifies what kind of action was audited.
type EventKind string

const (
	KindCheck       EventKind = "check"
	KindConfig      EventKind = "config"
	KindMaintenance EventKind = "maintenance"
	KindAlert       EventKind = "alert"
)

// Event represents a single auditable action.
type Event struct {
	ID        string    `json:"id"`
	Kind      EventKind `json:"kind"`
	Service   string    `json:"service"`
	Actor     string    `json:"actor"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Store holds an in-memory ring buffer of audit events.
type Store struct {
	mu     sync.RWMutex
	events []Event
	limit  int
	seq    int
}

// New returns a Store that retains at most limit events.
func New(limit int) *Store {
	if limit <= 0 {
		limit = 500
	}
	return &Store{limit: limit}
}

// Record appends an event to the store, evicting the oldest if at capacity.
func (s *Store) Record(kind EventKind, service, actor, message string) Event {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seq++
	e := Event{
		ID:        fmt.Sprintf("%d", s.seq),
		Kind:      kind,
		Service:   service,
		Actor:     actor,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	if len(s.events) >= s.limit {
		s.events = s.events[1:]
	}
	s.events = append(s.events, e)
	return e
}

// All returns a snapshot of all stored events, newest first.
func (s *Store) All() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, len(s.events))
	for i, e := range s.events {
		out[len(s.events)-1-i] = e
	}
	return out
}

// Filter returns events matching the given kind (empty string matches all).
func (s *Store) Filter(kind EventKind) []Event {
	all := s.All()
	if kind == "" {
		return all
	}
	var out []Event
	for _, e := range all {
		if e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}
