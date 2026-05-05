package scheduler

import (
	"encoding/json"
	"net/http"
	"time"
)

// StatusResponse describes the scheduler's runtime status.
type StatusResponse struct {
	Interval string `json:"interval"`
	Tracked  int    `json:"tracked_services"`
}

// NewStatusHandler returns an HTTP handler that exposes scheduler metadata.
func NewStatusHandler(s *Scheduler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		s.mu.Lock()
		tracked := len(s.prevStatus)
		s.mu.Unlock()

		resp := StatusResponse{
			Interval: s.interval.Round(time.Millisecond).String(),
			Tracked:  tracked,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
