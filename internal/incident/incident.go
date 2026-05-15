package incident

import (
	"errors"
	"sync"
	"time"
)

// Severity represents the impact level of an incident.
type Severity string

const (
	SeverityCritical Severity = "critical"
	SeverityHigh     Severity = "high"
	SeverityLow      Severity = "low"
)

// Incident represents a recorded service disruption.
type Incident struct {
	ID         string    `json:"id"`
	Service    string    `json:"service"`
	Title      string    `json:"title"`
	Severity   Severity  `json:"severity"`
	OpenedAt   time.Time `json:"opened_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	Resolved   bool      `json:"resolved"`
}

// Store holds all incidents in memory.
type Store struct {
	mu       sync.RWMutex
	records  map[string]*Incident
	counter  int
}

// New creates a new incident Store.
func New() *Store {
	return &Store{records: make(map[string]*Incident)}
}

// Open records a new incident and returns its ID.
func (s *Store) Open(service, title string, sev Severity) (string, error) {
	if service == "" {
		return "", errors.New("service is required")
	}
	if title == "" {
		return "", errors.New("title is required")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	id := fmt.Sprintf("INC-%04d", s.counter)
	s.records[id] = &Incident{
		ID:       id,
		Service:  service,
		Title:    title,
		Severity: sev,
		OpenedAt: time.Now().UTC(),
	}
	return id, nil
}

// Resolve marks an incident as resolved.
func (s *Store) Resolve(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	inc, ok := s.records[id]
	if !ok {
		return errors.New("incident not found")
	}
	if inc.Resolved {
		return errors.New("incident already resolved")
	}
	now := time.Now().UTC()
	inc.ResolvedAt = &now
	inc.Resolved = true
	return nil
}

// Get returns a single incident by ID.
func (s *Store) Get(id string) (*Incident, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	inc, ok := s.records[id]
	if !ok {
		return nil, false
	}
	copy := *inc
	return &copy, true
}

// Active returns all unresolved incidents.
func (s *Store) Active() []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*Incident
	for _, inc := range s.records {
		if !inc.Resolved {
			copy := *inc
			out = append(out, &copy)
		}
	}
	return out
}

// All returns every incident regardless of status.
func (s *Store) All() []*Incident {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Incident, 0, len(s.records))
	for _, inc := range s.records {
		copy := *inc
		out = append(out, &copy)
	}
	return out
}
