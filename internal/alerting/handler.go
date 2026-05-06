package alerting

import (
	"encoding/json"
	"net/http"
	"time"
)

type alertResponse struct {
	Service    string     `json:"service"`
	Message    string     `json:"message"`
	FiredAt    time.Time  `json:"fired_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

type handler struct {
	manager *Manager
}

// NewHandler returns an http.Handler that serves active alerts as JSON.
func NewHandler(m *Manager) http.Handler {
	return &handler{manager: m}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	active := h.manager.Active()
	resp := make([]alertResponse, 0, len(active))
	for _, a := range active {
		resp = append(resp, alertResponse{
			Service:    a.Service,
			Message:    a.Message,
			FiredAt:    a.FiredAt,
			ResolvedAt: a.ResolvedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
