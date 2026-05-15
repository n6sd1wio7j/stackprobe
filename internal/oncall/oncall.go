package oncall

import (
	"errors"
	"sync"
	"time"
)

// Schedule represents an on-call schedule for a service.
type Schedule struct {
	Service  string    `json:"service"`
	Owner    string    `json:"owner"`
	Phone    string    `json:"phone"`
	Email    string    `json:"email"`
	Until    time.Time `json:"until"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Store holds on-call schedules keyed by service name.
type Store struct {
	mu      sync.RWMutex
	records map[string]Schedule
}

// New returns an initialised Store.
func New() *Store {
	return &Store{records: make(map[string]Schedule)}
}

// Set registers or replaces the on-call schedule for a service.
func (s *Store) Set(svc string, owner, phone, email string, until time.Time) error {
	if svc == "" {
		return errors.New("service name required")
	}
	if owner == "" {
		return errors.New("owner required")
	}
	if until.IsZero() || until.Before(time.Now()) {
		return errors.New("until must be a future time")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[svc] = Schedule{
		Service:   svc,
		Owner:     owner,
		Phone:     phone,
		Email:     email,
		Until:     until,
		UpdatedAt: time.Now(),
	}
	return nil
}

// Get returns the on-call schedule for a service.
func (s *Store) Get(svc string) (Schedule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sc, ok := s.records[svc]
	return sc, ok
}

// Delete removes the on-call schedule for a service.
func (s *Store) Delete(svc string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.records[svc]; !ok {
		return errors.New("unknown service")
	}
	delete(s.records, svc)
	return nil
}

// All returns a snapshot of all on-call schedules.
func (s *Store) All() []Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Schedule, 0, len(s.records))
	for _, sc := range s.records {
		out = append(out, sc)
	}
	return out
}

// Active returns schedules whose Until time has not yet passed.
func (s *Store) Active() []Schedule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	var out []Schedule
	for _, sc := range s.records {
		if sc.Until.After(now) {
			out = append(out, sc)
		}
	}
	return out
}
