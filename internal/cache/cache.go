package cache

import (
	"sync"
	"time"
)

// Entry holds a cached check result with an expiry time.
type Entry struct {
	Status    string
	LatencyMs int64
	CachedAt  time.Time
	ExpiresAt time.Time
}

// IsExpired reports whether the cache entry is past its TTL.
func (e Entry) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// Cache stores the most recent health-check result per service
// and serves stale reads while a fresh check is in flight.
type Cache struct {
	mu      sync.RWMutex
	entries map[string]Entry
	ttl     time.Duration
}

// New returns a Cache with the given TTL applied to every entry.
func New(ttl time.Duration) *Cache {
	return &Cache{
		entries: make(map[string]Entry),
		ttl:     ttl,
	}
}

// Set stores a result for the named service, stamping it with the
// current time and computing an expiry from the configured TTL.
func (c *Cache) Set(service, status string, latencyMs int64) {
	now := time.Now()
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[service] = Entry{
		Status:    status,
		LatencyMs: latencyMs,
		CachedAt:  now,
		ExpiresAt: now.Add(c.ttl),
	}
}

// Get returns the entry for service and whether it exists.
// Callers should check Entry.IsExpired() to decide whether to
// perform a fresh probe.
func (c *Cache) Get(service string) (Entry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	e, ok := c.entries[service]
	return e, ok
}

// Delete removes the cached entry for service, forcing the next
// Get to report a miss.
func (c *Cache) Delete(service string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.entries, service)
}

// All returns a snapshot of every cached entry, keyed by service name.
func (c *Cache) All() map[string]Entry {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make(map[string]Entry, len(c.entries))
	for k, v := range c.entries {
		out[k] = v
	}
	return out
}
