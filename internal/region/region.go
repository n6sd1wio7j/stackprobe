package region

import (
	"errors"
	"sync"
)

// Region represents a named geographic or logical deployment zone.
type Region struct {
	Name     string   `json:"name"`
	Services []string `json:"services"`
}

// Store holds service-to-region mappings.
type Store struct {
	mu      sync.RWMutex
	regions map[string]string // service -> region name
}

// New creates an empty region Store.
func New() *Store {
	return &Store{
		regions: make(map[string]string),
	}
}

// Set assigns a region to a service.
func (s *Store) Set(service, region string) error {
	if service == "" {
		return errors.New("service name must not be empty")
	}
	if region == "" {
		return errors.New("region name must not be empty")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.regions[service] = region
	return nil
}

// Get returns the region assigned to a service.
func (s *Store) Get(service string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.regions[service]
	return r, ok
}

// Delete removes the region assignment for a service.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.regions, service)
}

// ByRegion returns all services grouped by region name.
func (s *Store) ByRegion() map[string][]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string][]string)
	for svc, reg := range s.regions {
		out[reg] = append(out[reg], svc)
	}
	return out
}

// Filter returns all services assigned to the given region.
func (s *Store) Filter(region string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []string
	for svc, reg := range s.regions {
		if reg == region {
			result = append(result, svc)
		}
	}
	return result
}
