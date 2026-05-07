package healthz

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// Response represents the health check response payload.
type Response struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version,omitempty"`
	GoVersion string    `json:"go_version"`
	Uptime    string    `json:"uptime"`
}

// Handler holds the state needed to serve liveness/readiness probes.
type Handler struct {
	version string
	started time.Time
}

// New creates a new Handler. version is the application version string.
func New(version string) *Handler {
	return &Handler{
		version: version,
		started: time.Now(),
	}
}

// LivenessHandler responds with 200 OK as long as the process is running.
func (h *Handler) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := Response{
		Status:    "ok",
		Timestamp: time.Now().UTC(),
		Version:   h.version,
		GoVersion: runtime.Version(),
		Uptime:    time.Since(h.started).Round(time.Second).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// ReadinessHandler responds with 200 OK when the service is ready to accept
// traffic. Callers may inject a readiness function to gate on dependencies.
func (h *Handler) ReadinessHandler(ready func() bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		status := "ok"
		code := http.StatusOK
		if ready != nil && !ready() {
			status = "unavailable"
			code = http.StatusServiceUnavailable
		}

		resp := Response{
			Status:    status,
			Timestamp: time.Now().UTC(),
			Version:   h.version,
			GoVersion: runtime.Version(),
			Uptime:    time.Since(h.started).Round(time.Second).String(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		_ = json.NewEncoder(w).Encode(resp)
	}
}
