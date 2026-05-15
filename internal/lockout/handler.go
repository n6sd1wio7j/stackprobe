package lockout

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type lockRequest struct {
	Reason   string `json:"reason"`
	Duration string `json:"duration"`
}

// NewHandler returns an http.Handler for managing service lockouts.
//
//	GET    /lockout          – list all active lockouts
//	PUT    /lockout/{service} – lock a service
//	DELETE /lockout/{service} – unlock a service
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/lockout", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/lockout/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/lockout/")
		if service == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPut:
			var req lockRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			d, err := time.ParseDuration(req.Duration)
			if err != nil || d <= 0 {
				http.Error(w, "invalid duration", http.StatusBadRequest)
				return
			}
			if err := s.Lock(service, req.Reason, d); err != nil {
				if err == ErrAlreadyLocked {
					http.Error(w, err.Error(), http.StatusConflict)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			if err := s.Unlock(service); err != nil {
				if err == ErrNotLocked {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
