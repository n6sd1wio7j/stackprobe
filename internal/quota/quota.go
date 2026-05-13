package quota

import (
	"errors"
	"sync"
	"time"
)

// ErrQuotaExceeded is returned when a service has exceeded its quota.
var ErrQuotaExceeded = errors.New("quota exceeded")

// ErrUnknownService is returned when no quota is configured for a service.
var ErrUnknownService = errors.New("unknown service")

// Entry holds quota configuration and current usage for a service.
type Entry struct {
	Service   string    `json:"service"`
	Limit     int       `json:"limit"`
	Used      int       `json:"used"`
	WindowEnd time.Time `json:"window_end"`
}

// Store manages per-service request quotas with a rolling window.
type Store struct {
	mu      sync.Mutex
	entries map[string]*Entry
	window  time.Duration
}

// New creates a Store with the given rolling window duration.
func New(window time.Duration) *Store {
	return &Store{
		entries: make(map[string]*Entry),
		window:  window,
	}
}

// Set configures a quota limit for the given service.
func (s *Store) Set(service string, limit int) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if limit <= 0 {
		return errors.New("limit must be greater than zero")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[service] = &Entry{
		Service:   service,
		Limit:     limit,
		Used:      0,
		WindowEnd: time.Now().Add(s.window),
	}
	return nil
}

// Allow checks whether the service may make another request.
// It resets the window if expired and increments usage on success.
func (s *Store) Allow(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[service]
	if !ok {
		return ErrUnknownService
	}
	now := time.Now()
	if now.After(e.WindowEnd) {
		e.Used = 0
		e.WindowEnd = now.Add(s.window)
	}
	if e.Used >= e.Limit {
		return ErrQuotaExceeded
	}
	e.Used++
	return nil
}

// Get returns the quota entry for the given service.
func (s *Store) Get(service string) (Entry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[service]
	if !ok {
		return Entry{}, ErrUnknownService
	}
	return *e, nil
}

// All returns a snapshot of all quota entries.
func (s *Store) All() []Entry {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]Entry, 0, len(s.entries))
	for _, e := range s.entries {
		out = append(out, *e)
	}
	return out
}
