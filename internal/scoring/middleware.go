package scoring

import (
	"net/http"
	"time"
)

// MetricsSource is satisfied by any type that can report uptime percentage
// and average latency for a named service.
type MetricsSource interface {
	UptimePercent(service string) float64
	AvgLatency(service string) time.Duration
}

// Refresher periodically recomputes scores from a MetricsSource.
type Refresher struct {
	store   *Store
	source  MetricsSource
	services []string
}

// NewRefresher creates a Refresher for the given services.
func NewRefresher(store *Store, source MetricsSource, services []string) *Refresher {
	return &Refresher{store: store, source: source, services: services}
}

// Refresh recomputes and stores the score for every registered service.
func (r *Refresher) Refresh() {
	for _, svc := range r.services {
		uptime := r.source.UptimePercent(svc)
		latencyMs := float64(r.source.AvgLatency(svc).Milliseconds())
		_ = r.store.Set(svc, uptime, latencyMs)
	}
}

// ScoreHeader is an HTTP middleware that injects an X-Service-Score header
// into the response based on the named service's current score.
func ScoreHeader(store *Store, service string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if sc, ok := store.Get(service); ok {
			w.Header().Set("X-Service-Score", sc.Grade)
		}
		next.ServeHTTP(w, r)
	})
}
