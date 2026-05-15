package badge

import (
	"errors"
	"fmt"
	"sync"
)

// Style represents the visual style of a badge.
type Style string

const (
	StyleFlat Style = "flat"
	StylePlastic Style = "plastic"
)

// Badge holds the metadata for a service status badge.
type Badge struct {
	Service string `json:"service"`
	Label   string `json:"label"`
	Status  string `json:"status"`
	Color   string `json:"color"`
	Style   Style  `json:"style"`
}

// Store manages badges for services.
type Store struct {
	mu     sync.RWMutex
	badges map[string]*Badge
}

// New creates a new badge Store.
func New() *Store {
	return &Store{badges: make(map[string]*Badge)}
}

// Set creates or updates the badge for a service.
func (s *Store) Set(service, label, status, color string, style Style) error {
	if service == "" {
		return errors.New("service must not be empty")
	}
	if label == "" {
		return errors.New("label must not be empty")
	}
	if status == "" {
		return errors.New("status must not be empty")
	}
	if color == "" {
		color = "grey"
	}
	if style == "" {
		style = StyleFlat
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.badges[service] = &Badge{
		Service: service,
		Label:   label,
		Status:  status,
		Color:   color,
		Style:   style,
	}
	return nil
}

// Get returns the badge for a service.
func (s *Store) Get(service string) (*Badge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.badges[service]
	if !ok {
		return nil, fmt.Errorf("no badge for service %q", service)
	}
	return b, nil
}

// Delete removes the badge for a service.
func (s *Store) Delete(service string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.badges[service]; !ok {
		return fmt.Errorf("no badge for service %q", service)
	}
	delete(s.badges, service)
	return nil
}

// All returns a snapshot of all badges.
func (s *Store) All() []*Badge {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Badge, 0, len(s.badges))
	for _, b := range s.badges {
		copy := *b
		out = append(out, &copy)
	}
	return out
}
