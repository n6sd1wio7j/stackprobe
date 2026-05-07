package statuspage

import (
	"sync"
	"time"
)

// Severity represents the impact level of an incident.
type Severity string

const (
	SeverityNone     Severity = "none"
	SeverityMinor    Severity = "minor"
	SeverityMajor    Severity = "major"
	SeverityCritical Severity = "critical"
)

// Incident represents a recorded service disruption or maintenance event.
type Incident struct {
	ID        string    `json:"id"`
	Service   string    `json:"service"`
	Title     string    `json:"title"`
	Severity  Severity  `json:"severity"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// Store holds active and resolved incidents.
type Store struct {
	mu        sync.RWMutex
	incidents map[string]*Incident
	nextID    int
}

// New creates an empty incident Store.
func New() *Store {
	return &Store{
		incidents: make(map[string]*Incident),
	}
}

// Add records a new incident and returns its assigned ID.
func (s *Store) Add(service, title, message string, severity Severity) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	id := fmt.Sprintf("INC-%04d", s.nextID)
	s.incidents[id] = &Incident{
		ID:        id,
		Service:   service,
		Title:     title,
		Severity:  severity,
		Message:   message,
		CreatedAt: time.Now().UTC(),
	}
	return id
}

// Resolve marks an incident as resolved by ID. Returns false if not found.
func (s *Store) Resolve(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	inc, ok := s.incidents[id]
	if !ok || inc.ResolvedAt != nil {
		return false
	}
	now := time.Now().UTC()
	inc.ResolvedAt = &now
	return true
}

// Active returns all unresolved incidents.
func (s *Store) Active() []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Incident
	for _, inc := range s.incidents {
		if inc.ResolvedAt == nil {
			out = append(out, inc)
		}
	}
	return out
}

// All returns every incident regardless of status.
func (s *Store) All() []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Incident, 0, len(s.incidents))
	for _, inc := range s.incidents {
		out = append(out, inc)
	}
	return out
}
