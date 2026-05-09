package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Record holds a point-in-time status snapshot for a single service.
type Record struct {
	Service   string        `json:"service"`
	Healthy   bool          `json:"healthy"`
	Latency   time.Duration `json:"latency_ns"`
	CapturedAt time.Time   `json:"captured_at"`
}

// Store persists and retrieves snapshots.
type Store struct {
	mu      sync.RWMutex
	latest  map[string]Record
}

// New returns an initialised Store.
func New() *Store {
	return &Store{latest: make(map[string]Record)}
}

// Save records the latest snapshot for a service, overwriting any previous one.
func (s *Store) Save(r Record) {
	if r.CapturedAt.IsZero() {
		r.CapturedAt = time.Now().UTC()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.latest[r.Service] = r
}

// Get returns the most recent snapshot for a service.
func (s *Store) Get(service string) (Record, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.latest[service]
	return r, ok
}

// All returns a copy of all stored snapshots.
func (s *Store) All() []Record {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Record, 0, len(s.latest))
	for _, r := range s.latest {
		out = append(out, r)
	}
	return out
}

// Dump writes all snapshots as JSON to the given file path.
func (s *Store) Dump(path string) error {
	records := s.All()
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
