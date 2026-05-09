// Package traceid provides request trace ID injection and propagation.
package traceid

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"sync"
)

const (
	Header        = "X-Trace-Id"
	contextKey    = contextKeyType("trace_id")
)

type contextKeyType string

// Store holds trace IDs associated with named services or requests.
type Store struct {
	mu  sync.RWMutex
	ids map[string]string
}

// New returns an initialised Store.
func New() *Store {
	return &Store{ids: make(map[string]string)}
}

// Generate creates a new random 16-byte hex trace ID.
func Generate() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Set stores a trace ID for the given key.
func (s *Store) Set(key, id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ids[key] = id
}

// Get returns the trace ID for the given key and whether it was found.
func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.ids[key]
	return v, ok
}

// Delete removes the trace ID for the given key.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ids, key)
}

// Middleware injects a trace ID into every request, reading from the incoming
// header first and generating one if absent.
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(Header)
		if id == "" {
			var err error
			id, err = Generate()
			if err != nil {
				id = "unknown"
			}
		}
		w.Header().Set(Header, id)
		ctx := context.WithValue(r.Context(), contextKey, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// FromContext extracts the trace ID from a context, returning empty string if absent.
func FromContext(ctx context.Context) string {
	v, _ := ctx.Value(contextKey).(string)
	return v
}
