package mute

import (
	"encoding/json"
	"net/http"
	"time"
)

// NewHandler returns an http.Handler that exposes mute management endpoints.
//
//	GET    /mute          – list all active mute windows
//	PUT    /mute          – create or replace a mute window
//	DELETE /mute?service= – remove a mute window
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/mute", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, s.All(), http.StatusOK)
		case http.MethodPut:
			var body struct {
				Service string `json:"service"`
				Reason  string `json:"reason"`
				Until   string `json:"until"` // RFC3339
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			until, err := time.Parse(time.RFC3339, body.Until)
			if err != nil {
				http.Error(w, "invalid until: use RFC3339", http.StatusBadRequest)
				return
			}
			if err := s.Mute(body.Service, body.Reason, until); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			svc := r.URL.Query().Get("service")
			if svc == "" {
				http.Error(w, "service query param required", http.StatusBadRequest)
				return
			}
			if err := s.Unmute(svc); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
