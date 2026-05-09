package audit

import (
	"encoding/json"
	"net/http"
)

// NewHandler returns an http.Handler that exposes the audit log over HTTP.
//
// GET /audit          — returns all events (newest first)
// GET /audit?kind=... — filters by event kind
func NewHandler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		kind := EventKind(r.URL.Query().Get("kind"))
		events := s.Filter(kind)
		if events == nil {
			events = []Event{}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(events)
	})
}
