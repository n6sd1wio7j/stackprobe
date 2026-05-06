package retry

import (
	"encoding/json"
	"net/http"
)

// policyResponse is the JSON representation of a retry policy.
type policyResponse struct {
	MaxAttempts int     `json:"max_attempts"`
	DelayMs     int64   `json:"delay_ms"`
	Backoff     float64 `json:"backoff"`
}

// NewHandler returns an HTTP handler that exposes the active retry policy.
func NewHandler(r *Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		resp := policyResponse{
			MaxAttempts: r.policy.MaxAttempts,
			DelayMs:     r.policy.Delay.Milliseconds(),
			Backoff:     r.policy.Backoff,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	}
}
