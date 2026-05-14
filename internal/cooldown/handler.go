package cooldown

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type handler struct {
	store *Store
}

// NewHandler returns an http.Handler that exposes cooldown state over HTTP.
//
//	GET  /cooldown          – list all entries
//	PUT  /cooldown/{svc}    – set custom duration (body: {"duration":"30s"})
//	DELETE /cooldown/{svc}  – reset cooldown for service
func NewHandler(s *Store) http.Handler {
	h := &handler{store: s}
	mux := http.NewServeMux()
	mux.HandleFunc("/cooldown", h.list)
	mux.HandleFunc("/cooldown/", h.dispatch)
	return mux
}

func (h *handler) list(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	h.store.mu.Lock()
	type row struct {
		Service   string  `json:"service"`
		Duration  string  `json:"duration"`
		Remaining string  `json:"remaining"`
	}
	rows := make([]row, 0, len(h.store.entries))
	for svc, e := range h.store.entries {
		elapsed := time.Since(e.lastTriggered)
		rem := time.Duration(0)
		if !e.lastTriggered.IsZero() && elapsed < e.duration {
			rem = e.duration - elapsed
		}
		rows = append(rows, row{
			Service:   svc,
			Duration:  e.duration.String(),
			Remaining: rem.String(),
		})
	}
	h.store.mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rows)
}

func (h *handler) dispatch(w http.ResponseWriter, r *http.Request) {
	svc := strings.TrimPrefix(r.URL.Path, "/cooldown/")
	if svc == "" {
		http.Error(w, "service name required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var body struct {
			Duration string `json:"duration"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		d, err := time.ParseDuration(body.Duration)
		if err != nil {
			http.Error(w, "invalid duration", http.StatusBadRequest)
			return
		}
		if err := h.store.Set(svc, d); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if err := h.store.Reset(svc); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
