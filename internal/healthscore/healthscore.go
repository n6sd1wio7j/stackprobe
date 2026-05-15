package healthscore

import (
	"errors"
	"sync"
	"time"
)

// Score holds the computed health score for a service.
type Score struct {
	Service   string
	Value     float64 // 0.0 – 100.0
	Grade     string  // A, B, C, D, F
	ComputedAt time.Time
}

// Store computes and caches health scores derived from uptime and latency.
type Store struct {
	mu     sync.RWMutex
	scores map[string]Score
}

// New returns an initialised Store.
func New() *Store {
	return &Store{scores: make(map[string]Score)}
}

// Record computes a score from uptimePct (0–100) and avgLatencyMs, then
// stores it for the given service.
func (s *Store) Record(service string, uptimePct float64, avgLatencyMs float64) error {
	if service == "" {
		return errors.New("healthscore: service name required")
	}
	if uptimePct < 0 || uptimePct > 100 {
		return errors.New("healthscore: uptimePct must be in [0, 100]")
	}

	// Latency penalty: subtract up to 20 points for latency >= 2000 ms.
	latencyPenalty := avgLatencyMs / 2000.0 * 20.0
	if latencyPenalty > 20 {
		latencyPenalty = 20
	}

	value := uptimePct - latencyPenalty
	if value < 0 {
		value = 0
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[service] = Score{
		Service:    service,
		Value:      value,
		Grade:      grade(value),
		ComputedAt: time.Now().UTC(),
	}
	return nil
}

// Get returns the stored Score for a service.
func (s *Store) Get(service string) (Score, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sc, ok := s.scores[service]
	return sc, ok
}

// All returns a snapshot of every stored score.
func (s *Store) All() []Score {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Score, 0, len(s.scores))
	for _, sc := range s.scores {
		out = append(out, sc)
	}
	return out
}

func grade(v float64) string {
	switch {
	case v >= 95:
		return "A"
	case v >= 85:
		return "B"
	case v >= 70:
		return "C"
	case v >= 50:
		return "D"
	default:
		return "F"
	}
}
