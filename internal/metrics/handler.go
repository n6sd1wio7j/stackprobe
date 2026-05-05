package metrics

import (
	"encoding/json"
	"net/http"
	"sort"
)

type metricsResponse struct {
	Service      string  `json:"service"`
	TotalChecks  int64   `json:"total_checks"`
	UpCount      int64   `json:"up_count"`
	DownCount    int64   `json:"down_count"`
	UptimePct    float64 `json:"uptime_percent"`
	AvgLatencyMs float64 `json:"avg_latency_ms"`
	LastChecked  string  `json:"last_checked"`
}

// NewHandler returns an HTTP handler that serves metrics as JSON.
func NewHandler(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		all := store.All()
		sort.Slice(all, func(i, j int) bool {
			return all[i].ServiceName < all[j].ServiceName
		})

		responses := make([]metricsResponse, 0, len(all))
		for _, m := range all {
			responses = append(responses, metricsResponse{
				Service:      m.ServiceName,
				TotalChecks:  m.TotalChecks,
				UpCount:      m.UpCount,
				DownCount:    m.DownCount,
				UptimePct:    m.UptimePercent(),
				AvgLatencyMs: float64(m.AvgLatency().Milliseconds()),
				LastChecked:  m.LastChecked.UTC().Format("2006-01-02T15:04:05Z"),
			})
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(responses); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
