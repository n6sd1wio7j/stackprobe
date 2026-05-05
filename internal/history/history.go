package history

import (
	"sync"
	"time"
)

// Record holds the result of a single health check at a point in time.
type Record struct {
	Service   string    `json:"service"`
	Healthy   bool      `json:"healthy"`
	Status    int       `json:"status"`
	CheckedAt time.Time `json:"checked_at"`
}

// Store keeps a bounded in-memory history of health check records per service.
type Store struct {
	mu      sync.RWMutex
	records map[string][]Record
	limit   int
}

// New creates a new Store that retains at most limit records per service.
func New(limit int) *Store {
	if limit <= 0 {
		limit = 10
	}
	return &Store{
		records: make(map[string][]Record),
		limit:   limit,
	}
}

// Add appends a new record for the given service, evicting the oldest if the
// limit has been reached.
func (s *Store) Add(r Record) {
	s.mu.Lock()
	defer s.mu.Unlock()

	bucket := s.records[r.Service]
	bucket = append(bucket, r)
	if len(bucket) > s.limit {
		bucket = bucket[len(bucket)-s.limit:]
	}
	s.records[r.Service] = bucket
}

// Get returns a copy of all records stored for the given service.
func (s *Store) Get(service string) []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	src := s.records[service]
	out := make([]Record, len(src))
	copy(out, src)
	return out
}

// All returns a copy of all records across every service.
func (s *Store) All() map[string][]Record {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make(map[string][]Record, len(s.records))
	for svc, recs := range s.records {
		bucket := make([]Record, len(recs))
		copy(bucket, recs)
		out[svc] = bucket
	}
	return out
}
