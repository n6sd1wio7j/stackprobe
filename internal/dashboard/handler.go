package dashboard

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/stackprobe/internal/aggregator"
)

// Status represents the overall dashboard response payload.
type Status struct {
	Timestamp time.Time                   `json:"timestamp"`
	Healthy   int                         `json:"healthy"`
	Unhealthy int                         `json:"unhealthy"`
	Services  []aggregator.ServiceStatus  `json:"services"`
}

// Collector is the interface satisfied by aggregator.Aggregator.
type Collector interface {
	Collect() []aggregator.ServiceStatus
}

// Handler holds dependencies for the dashboard HTTP handler.
type Handler struct {
	collector Collector
}

// New creates a new Handler with the given Collector.
func New(c Collector) *Handler {
	return &Handler{collector: c}
}

// ServeHTTP handles GET /status requests and returns aggregated service health.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	services := h.collector.Collect()

	var healthy, unhealthy int
	for _, s := range services {
		if s.Up {
			healthy++
		} else {
			unhealthy++
		}
	}

	payload := Status{
		Timestamp: time.Now().UTC(),
		Healthy:   healthy,
		Unhealthy: unhealthy,
		Services:  services,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(payload)
}
