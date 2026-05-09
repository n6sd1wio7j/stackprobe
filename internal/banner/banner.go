// Package banner provides service announcement and message-of-the-day
// functionality for the stackprobe dashboard.
package banner

import (
	"sync"
	"time"
)

// Banner holds an announcement message for the dashboard.
type Banner struct {
	ID        string    `json:"id"`
	Message   string    `json:"message"`
	Level     string    `json:"level"` // info, warning, critical
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
}

// Store manages active banners.
type Store struct {
	mu      sync.RWMutex
	banners map[string]*Banner
	nextID  int
}

// New creates a new banner Store.
func New() *Store {
	return &Store{
		banners: make(map[string]*Banner),
	}
}

// Add creates a new banner and stores it, returning the assigned ID.
func (s *Store) Add(message, level string, expiresAt time.Time) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	id := fmt.Sprintf("banner-%d", s.nextID)
	s.banners[id] = &Banner{
		ID:        id,
		Message:   message,
		Level:     level,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
	}
	return id
}

// Remove deletes a banner by ID. Returns false if not found.
func (s *Store) Remove(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.banners[id]
	if ok {
		delete(s.banners, id)
	}
	return ok
}

// Active returns all non-expired banners.
func (s *Store) Active() []*Banner {
	s.mu.RLock()
	defer s.mu.RUnlock()

	now := time.Now()
	var result []*Banner
	for _, b := range s.banners {
		if b.ExpiresAt.IsZero() || b.ExpiresAt.After(now) {
			result = append(result, b)
		}
	}
	return result
}
