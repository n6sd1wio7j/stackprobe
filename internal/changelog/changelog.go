package changelog

import (
	"errors"
	"sync"
	"time"
)

// EntryKind classifies a changelog entry.
type EntryKind string

const (
	KindAdded    EntryKind = "added"
	KindChanged  EntryKind = "changed"
	KindFixed    EntryKind = "fixed"
	KindRemoved  EntryKind = "removed"
	KindDegraded EntryKind = "degraded"
)

// Entry represents a single changelog event for a service.
type Entry struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Kind      EntryKind `json:"kind"`
	Message   string    `json:"message"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"created_at"`
}

// Store holds changelog entries per service.
type Store struct {
	mu      sync.RWMutex
	entries []Entry
	limit   int
	nextID  int
}

// New creates a Store that retains at most limit entries.
func New(limit int) *Store {
	if limit <= 0 {
		limit = 200
	}
	return &Store{limit: limit}
}

// Add appends a new entry for the given service.
func (s *Store) Add(service string, kind EntryKind, message, author string) (Entry, error) {
	if service == "" {
		return Entry{}, errors.New("service is required")
	}
	if message == "" {
		return Entry{}, errors.New("message is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	e := Entry{
		ID:        fmt.Sprintf("%d", s.nextID),
		Service:   service,
		Kind:      kind,
		Message:   message,
		Author:    author,
		CreatedAt: time.Now().UTC(),
	}
	s.entries = append([]Entry{e}, s.entries...)
	if len(s.entries) > s.limit {
		s.entries = s.entries[:s.limit]
	}
	return e, nil
}

// Get returns all entries for a service, newest first.
func (s *Store) Get(service string) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Entry
	for _, e := range s.entries {
		if e.Service == service {
			out = append(out, e)
		}
	}
	return out
}

// All returns every entry, newest first.
func (s *Store) All() []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}

// Filter returns entries whose Kind matches kind. If kind is empty, all are returned.
func (s *Store) Filter(kind EntryKind) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Entry
	for _, e := range s.entries {
		if kind == "" || e.Kind == kind {
			out = append(out, e)
		}
	}
	return out
}
