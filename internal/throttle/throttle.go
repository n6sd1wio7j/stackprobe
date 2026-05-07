package throttle

import (
	"net/http"
	"sync"
	"time"
)

// Store tracks per-service request throttle state.
type Store struct {
	mu       sync.Mutex
	entries  map[string]*entry
	default_ int // max requests per window
	window   time.Duration
}

type entry struct {
	count     int
	windowEnd time.Time
}

// New creates a Store with the given default limit and window duration.
func New(defaultLimit int, window time.Duration) *Store {
	return &Store{
		entries:  make(map[string]*entry),
		default_: defaultLimit,
		window:   window,
	}
}

// Allow reports whether a request for the given service is permitted under
// the throttle limit. It increments the counter and returns false when the
// limit has been reached for the current window.
func (s *Store) Allow(service string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	e, ok := s.entries[service]
	if !ok || now.After(e.windowEnd) {
		s.entries[service] = &entry{count: 1, windowEnd: now.Add(s.window)}
		return true
	}
	if e.count >= s.default_ {
		return false
	}
	e.count++
	return true
}

// Remaining returns how many requests are left in the current window.
func (s *Store) Remaining(service string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	e, ok := s.entries[service]
	if !ok || now.After(e.windowEnd) {
		return s.default_
	}
	remaining := s.default_ - e.count
	if remaining < 0 {
		return 0
	}
	return remaining
}

// Reset clears the throttle state for a service, allowing immediate requests.
func (s *Store) Reset(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, service)
}

// Middleware returns an HTTP middleware that throttles requests by the
// "X-Service" header value. Returns 429 when the limit is exceeded.
func (s *Store) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		service := r.Header.Get("X-Service")
		if service == "" {
			service = "default"
		}
		if !s.Allow(service) {
			http.Error(w, "throttle limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
