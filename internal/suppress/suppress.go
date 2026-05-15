package suppress

import (
	"errors"
	"sync"
	"time"
)

// Rule holds a suppression rule for a service.
type Rule struct {
	Service   string
	Reason    string
	ExpiresAt time.Time
	CreatedAt time.Time
}

// IsExpired reports whether the suppression window has passed.
func (r Rule) IsExpired(now time.Time) bool {
	return now.After(r.ExpiresAt)
}

// Store manages alert-suppression windows per service.
type Store struct {
	mu    sync.RWMutex
	rules map[string]Rule
}

// New returns an initialised Store.
func New() *Store {
	return &Store{rules: make(map[string]Rule)}
}

// Suppress adds or replaces a suppression window for service.
func (s *Store) Suppress(service, reason string, until time.Time) error {
	if service == "" {
		return errors.New("service name required")
	}
	if reason == "" {
		return errors.New("reason required")
	}
	if !until.After(time.Now()) {
		return errors.New("expiry must be in the future")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[service] = Rule{
		Service:   service,
		Reason:    reason,
		ExpiresAt: until,
		CreatedAt: time.Now(),
	}
	return nil
}

// Unsuppress removes a suppression window for service.
func (s *Store) Unsuppress(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[service]; !ok {
		return errors.New("no suppression found for service")
	}
	delete(s.rules, service)
	return nil
}

// IsSuppressed reports whether service has an active (non-expired) window.
func (s *Store) IsSuppressed(service string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[service]
	if !ok {
		return false
	}
	return !r.IsExpired(time.Now())
}

// Get returns the rule for service and whether it exists.
func (s *Store) Get(service string) (Rule, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[service]
	return r, ok
}

// All returns a snapshot of every stored rule (including expired).
func (s *Store) All() []Rule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Rule, 0, len(s.rules))
	for _, r := range s.rules {
		out = append(out, r)
	}
	return out
}
