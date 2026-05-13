package heartbeat

import (
	"encoding/json"
	"net/http"
	"strings"
)

type handler struct {
	store *Store
}

// NewHandler returns an http.Handler for the heartbeat API.
// POST /heartbeat/{service}  — record a ping
// GET  /heartbeat            — list all beats
// GET  /heartbeat/stale      — list stale beats
func NewHandler(s *Store) http.Handler {
	return &handler{store: s}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/heartbeat")
	path = strings.TrimPrefix(path, "/")

	switch r.Method {
	case http.MethodPost:
		if path == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		h.store.Ping(path)
		w.WriteHeader(http.StatusNoContent)

	case http.MethodGet:
		if path == "stale" {
			h.writeJSON(w, h.store.Stale())
			return
		}
		h.writeJSON(w, h.store.All())

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
