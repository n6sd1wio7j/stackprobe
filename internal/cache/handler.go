package cache

import (
	"encoding/json"
	"net/http"
)

type entryResponse struct {
	Service   string `json:"service"`
	Status    string `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	CachedAt  string `json:"cached_at"`
	ExpiresAt string `json:"expires_at"`
	Expired   bool   `json:"expired"`
}

// NewHandler returns an HTTP handler that exposes the current cache
// snapshot as JSON. GET /cache lists every entry; individual service
// look-ups are done via ?service=<name>.
func NewHandler(c *Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var entries []entryResponse

		if svc := r.URL.Query().Get("service"); svc != "" {
			e, ok := c.Get(svc)
			if !ok {
				http.Error(w, "service not found", http.StatusNotFound)
				return
			}
			entries = []entryResponse{toResponse(svc, e)}
		} else {
			all := c.All()
			entries = make([]entryResponse, 0, len(all))
			for svc, e := range all {
				entries = append(entries, toResponse(svc, e))
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(entries) //nolint:errcheck
	}
}

func toResponse(service string, e Entry) entryResponse {
	return entryResponse{
		Service:   service,
		Status:    e.Status,
		LatencyMs: e.LatencyMs,
		CachedAt:  e.CachedAt.UTC().Format("2006-01-02T15:04:05Z"),
		ExpiresAt: e.ExpiresAt.UTC().Format("2006-01-02T15:04:05Z"),
		Expired:   e.IsExpired(),
	}
}
