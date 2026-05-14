package scoring

import (
	"errors"
	"sync"
)

// Score holds the computed health score for a service.
type Score struct {
	Service  string  `json:"service"`
	Value    float64 `json:"value"`    // 0.0 – 100.0
	Grade    string  `json:"grade"`    // A, B, C, D, F
	Reasoning string `json:"reasoning"`
}

// Store holds per-service scores.
type Store struct {
	mu     sync.RWMutex
	scores map[string]Score
}

// New returns an initialised Store.
func New() *Store {
	return &Store{scores: make(map[string]Score)}
}

// Set computes and stores the score for a service given uptime (0-100)
// and average latency in milliseconds.
func (s *Store) Set(service string, uptimePct float64, avgLatencyMs float64) error {
	if service == "" {
		return errors.New("service name required")
	}
	if uptimePct < 0 || uptimePct > 100 {
		return errors.New("uptimePct must be between 0 and 100")
	}

	// Weighted formula: 70 % uptime + 30 % latency score.
	latencyScore := latencyPoints(avgLatencyMs)
	value := 0.70*uptimePct + 0.30*latencyScore
	if value > 100 {
		value = 100
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.scores[service] = Score{
		Service:   service,
		Value:     value,
		Grade:     grade(value),
		Reasoning: reasoning(uptimePct, avgLatencyMs),
	}
	return nil
}

// Get returns the score for a service.
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

// Delete removes the score for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.scores[service]; !ok {
		return errors.New("unknown service")
	}
	delete(s.scores, service)
	return nil
}

// latencyPoints converts avg latency (ms) into a 0-100 score (lower is better).
func latencyPoints(ms float64) float64 {
	switch {
	case ms <= 50:
		return 100
	case ms <= 200:
		return 80
	case ms <= 500:
		return 60
	case ms <= 1000:
		return 40
	default:
		return 10
	}
}

func grade(v float64) string {
	switch {
	case v >= 90:
		return "A"
	case v >= 75:
		return "B"
	case v >= 60:
		return "C"
	case v >= 40:
		return "D"
	default:
		return "F"
	}
}

func reasoning(uptime, latencyMs float64) string {
	if uptime < 90 {
		return "low uptime dragging score down"
	}
	if latencyMs > 500 {
		return "high latency impacting score"
	}
	return "healthy"
}
