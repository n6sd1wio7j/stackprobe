package uptime

import (
	"sync"
	"time"
)

// Record holds uptime tracking data for a single service.
type Record struct {
	Service   string
	StartedAt time.Time
	UpSince   time.Time
	DownSince *time.Time
	IsUp      bool
}

// Tracker maintains uptime state for multiple services.
type Tracker struct {
	mu      sync.RWMutex
	records map[string]*Record
}

// New creates a new Tracker.
func New() *Tracker {
	return &Tracker{
		records: make(map[string]*Record),
	}
}

// Report updates the uptime state for the given service.
func (t *Tracker) Report(service string, up bool) {
	t.mu.Lock()
	defer t.mu.Unlock()

	now := time.Now()
	r, exists := t.records[service]
	if !exists {
		r = &Record{
			Service:   service,
			StartedAt: now,
			UpSince:   now,
			IsUp:      up,
		}
		if !up {
			r.DownSince = &now
		}
		t.records[service] = r
		return
	}

	if up && !r.IsUp {
		r.IsUp = true
		r.UpSince = now
		r.DownSince = nil
	} else if !up && r.IsUp {
		r.IsUp = false
		r.DownSince = &now
	}
}

// Get returns the uptime record for a service, or false if unknown.
func (t *Tracker) Get(service string) (Record, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.records[service]
	if !ok {
		return Record{}, false
	}
	return *r, true
}

// All returns a snapshot of all tracked services.
func (t *Tracker) All() []Record {
	t.mu.RLock()
	defer t.mu.RUnlock()
	out := make([]Record, 0, len(t.records))
	for _, r := range t.records {
		out = append(out, *r)
	}
	return out
}

// Duration returns how long the service has been in its current state.
func (t *Tracker) Duration(service string) (time.Duration, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	r, ok := t.records[service]
	if !ok {
		return 0, false
	}
	if r.IsUp {
		return time.Since(r.UpSince), true
	}
	return time.Since(*r.DownSince), true
}
