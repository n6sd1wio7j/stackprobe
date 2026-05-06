package ratelimit

import (
	"sync"
	"time"
)

// Limiter enforces a minimum interval between successive checks for each service.
type Limiter struct {
	mu       sync.Mutex
	last     map[string]time.Time
	interval time.Duration
}

// New creates a Limiter that allows at most one check per interval per service.
func New(interval time.Duration) *Limiter {
	if interval <= 0 {
		interval = time.Second
	}
	return &Limiter{
		last:     make(map[string]time.Time),
		interval: interval,
	}
}

// Allow reports whether the named service may be checked right now.
// If allowed, it records the current time as the last check time.
func (l *Limiter) Allow(service string) bool {
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()

	if t, ok := l.last[service]; ok && now.Sub(t) < l.interval {
		return false
	}
	l.last[service] = now
	return true
}

// Reset clears the last-seen time for the named service,
// allowing it to be checked immediately on the next call.
func (l *Limiter) Reset(service string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, service)
}

// NextAllowed returns the earliest time at which the named service
// may be checked. Returns the zero time if the service has never been checked.
func (l *Limiter) NextAllowed(service string) time.Time {
	l.mu.Lock()
	defer l.mu.Unlock()

	if t, ok := l.last[service]; ok {
		return t.Add(l.interval)
	}
	return time.Time{}
}
