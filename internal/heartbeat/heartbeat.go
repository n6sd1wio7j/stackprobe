package heartbeat

import (
	"sync"
	"time"
)

// Beat records the last seen time for a service.
type Beat struct {
	Service   string
	LastSeen  time.Time
	Interval  time.Duration
}

// IsStale returns true if the beat has not been updated within its interval.
func (b Beat) IsStale(now time.Time) bool {
	return now.Sub(b.LastSeen) > b.Interval
}

// Store holds heartbeat records for registered services.
type Store struct {
	mu      sync.RWMutex
	beats   map[string]Beat
	default_ time.Duration
}

// New creates a Store with the given default stale interval.
func New(defaultInterval time.Duration) *Store {
	return &Store{
		beats:    make(map[string]Beat),
		default_: defaultInterval,
	}
}

// Ping records the current time as the last heartbeat for service.
func (s *Store) Ping(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.beats[service]
	interval := s.default_
	if ok {
		interval = existing.Interval
	}
	s.beats[service] = Beat{
		Service:  service,
		LastSeen: time.Now(),
		Interval: interval,
	}
}

// Register adds a service with a custom stale interval without setting LastSeen.
func (s *Store) Register(service string, interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	b := s.beats[service]
	b.Service = service
	b.Interval = interval
	s.beats[service] = b
}

// Get returns the Beat for a service and whether it exists.
func (s *Store) Get(service string) (Beat, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.beats[service]
	return b, ok
}

// Stale returns all services whose heartbeat has gone stale.
func (s *Store) Stale() []Beat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now()
	var out []Beat
	for _, b := range s.beats {
		if !b.LastSeen.IsZero() && b.IsStale(now) {
			out = append(out, b)
		}
	}
	return out
}

// All returns every registered Beat.
func (s *Store) All() []Beat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Beat, 0, len(s.beats))
	for _, b := range s.beats {
		out = append(out, b)
	}
	return out
}
