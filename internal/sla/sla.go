package sla

import (
	"sync"
	"time"
)

// Target holds the SLA configuration for a single service.
type Target struct {
	Service       string
	UptimePercent float64 // e.g. 99.9
	MaxLatencyMs  int64   // 0 means no latency SLA
}

// Violation records a breach of an SLA target.
type Violation struct {
	Service   string
	Kind      string // "uptime" or "latency"
	Message   string
	OccuredAt time.Time
}

// Store manages SLA targets and accumulated violations.
type Store struct {
	mu         sync.RWMutex
	targets    map[string]Target
	violations []Violation
}

// New creates an empty SLA store.
func New() *Store {
	return &Store{
		targets: make(map[string]Target),
	}
}

// Set registers or replaces the SLA target for a service.
func (s *Store) Set(t Target) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.targets[t.Service] = t
}

// Get returns the target for a service and whether it exists.
func (s *Store) Get(service string) (Target, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.targets[service]
	return t, ok
}

// Evaluate checks current uptime and latency against the target and records
// any violations. It is a no-op if no target is registered for the service.
func (s *Store) Evaluate(service string, uptimePct float64, avgLatencyMs int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.targets[service]
	if !ok {
		return
	}

	now := time.Now().UTC()

	if uptimePct < t.UptimePercent {
		s.violations = append(s.violations, Violation{
			Service:   service,
			Kind:      "uptime",
			Message:   "uptime below target",
			OccuredAt: now,
		})
	}

	if t.MaxLatencyMs > 0 && avgLatencyMs > t.MaxLatencyMs {
		s.violations = append(s.violations, Violation{
			Service:   service,
			Kind:      "latency",
			Message:   "average latency exceeds target",
			OccuredAt: now,
		})
	}
}

// Violations returns a snapshot of all recorded violations.
func (s *Store) Violations() []Violation {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Violation, len(s.violations))
	copy(out, s.violations)
	return out
}

// Targets returns a snapshot of all registered targets.
func (s *Store) Targets() []Target {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Target, 0, len(s.targets))
	for _, t := range s.targets {
		out = append(out, t)
	}
	return out
}
