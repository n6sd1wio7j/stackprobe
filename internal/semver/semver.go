// Package semver tracks and compares service version strings.
package semver

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
)

// Version holds a parsed semantic version.
type Version struct {
	Major, Minor, Patch int
	Raw                 string
}

// Parse parses a semver string like "1.2.3" or "v1.2.3".
func Parse(s string) (Version, error) {
	trimmed := strings.TrimPrefix(s, "v")
	parts := strings.SplitN(trimmed, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("semver: invalid version %q", s)
	}
	var v Version
	v.Raw = s
	var err error
	if v.Major, err = strconv.Atoi(parts[0]); err != nil {
		return Version{}, fmt.Errorf("semver: invalid major in %q", s)
	}
	if v.Minor, err = strconv.Atoi(parts[1]); err != nil {
		return Version{}, fmt.Errorf("semver: invalid minor in %q", s)
	}
	if v.Patch, err = strconv.Atoi(parts[2]); err != nil {
		return Version{}, fmt.Errorf("semver: invalid patch in %q", s)
	}
	return v, nil
}

// Compare returns -1, 0, or 1 if a is less than, equal to, or greater than b.
func Compare(a, b Version) int {
	for _, pair := range [][2]int{{a.Major, b.Major}, {a.Minor, b.Minor}, {a.Patch, b.Patch}} {
		if pair[0] < pair[1] {
			return -1
		}
		if pair[0] > pair[1] {
			return 1
		}
	}
	return 0
}

// Store maps service names to their reported version.
type Store struct {
	mu   sync.RWMutex
	data map[string]Version
}

// New returns an initialised Store.
func New() *Store {
	return &Store{data: make(map[string]Version)}
}

var errEmptyService = errors.New("semver: service name must not be empty")

// Set records the version string for a service.
func (s *Store) Set(service, raw string) error {
	if service == "" {
		return errEmptyService
	}
	v, err := Parse(raw)
	if err != nil {
		return err
	}
	s.mu.Lock()
	s.data[service] = v
	s.mu.Unlock()
	return nil
}

// Get returns the Version for the given service.
func (s *Store) Get(service string) (Version, bool) {
	s.mu.RLock()
	v, ok := s.data[service]
	s.mu.RUnlock()
	return v, ok
}

// Delete removes a service entry.
func (s *Store) Delete(service string) {
	s.mu.Lock()
	delete(s.data, service)
	s.mu.Unlock()
}

// All returns a snapshot of all service versions.
func (s *Store) All() map[string]Version {
	s.mu.RLock()
	out := make(map[string]Version, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	s.mu.RUnlock()
	return out
}
