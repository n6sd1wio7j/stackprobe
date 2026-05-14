package deadletter

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the dead-letter store.
//
//	GET  /deadletter          — all entries across services
//	GET  /deadletter/{service} — entries for a specific service
//	DELETE /deadletter/{service} — flush entries for a service
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/deadletter", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/deadletter/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/deadletter/")
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			ents, err := s.Get(svc)
			if err == ErrUnknownService {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, ents)
		case http.MethodDelete:
			ents, err := s.Flush(svc)
			if err == ErrUnknownService {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, ents)
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
