package metrics

import (
	"sync"
	"time"
)

// ServiceMetrics holds aggregated metrics for a single service.
type ServiceMetrics struct {
	ServiceName  string
	TotalChecks  int64
	UpCount      int64
	DownCount    int64
	TotalLatency time.Duration
	LastChecked  time.Time
}

// UptimePercent returns the percentage of successful checks.
func (m *ServiceMetrics) UptimePercent() float64 {
	if m.TotalChecks == 0 {
		return 0
	}
	return float64(m.UpCount) / float64(m.TotalChecks) * 100
}

// AvgLatency returns the average latency across all checks.
func (m *ServiceMetrics) AvgLatency() time.Duration {
	if m.TotalChecks == 0 {
		return 0
	}
	return time.Duration(int64(m.TotalLatency) / m.TotalChecks)
}

// Store holds metrics for all tracked services.
type Store struct {
	mu      sync.RWMutex
	records map[string]*ServiceMetrics
}

// New creates a new metrics Store.
func New() *Store {
	return &Store{
		records: make(map[string]*ServiceMetrics),
	}
}

// Record updates metrics for a service after a health check.
func (s *Store) Record(name string, up bool, latency time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	m, ok := s.records[name]
	if !ok {
		m = &ServiceMetrics{ServiceName: name}
		s.records[name] = m
	}

	m.TotalChecks++
	m.TotalLatency += latency
	m.LastChecked = time.Now()
	if up {
		m.UpCount++
	} else {
		m.DownCount++
	}
}

// Get returns a copy of metrics for a named service.
func (s *Store) Get(name string) (ServiceMetrics, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.records[name]
	if !ok {
		return ServiceMetrics{}, false
	}
	return *m, true
}

// All returns a copy of all service metrics.
func (s *Store) All() []ServiceMetrics {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]ServiceMetrics, 0, len(s.records))
	for _, m := range s.records {
		out = append(out, *m)
	}
	return out
}
