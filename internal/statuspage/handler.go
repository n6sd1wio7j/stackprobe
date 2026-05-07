package statuspage

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler for the status-page incident API.
//
// Routes:
//   GET  /incidents        — list all incidents
//   GET  /incidents/active — list active (unresolved) incidents
//   POST /incidents        — create a new incident
//   POST /incidents/{id}/resolve — resolve an incident
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/incidents", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, s.All())
		case http.MethodPost:
			var body struct {
				Service  string   `json:"service"`
				Title    string   `json:"title"`
				Message  string   `json:"message"`
				Severity Severity `json:"severity"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if body.Service == "" || body.Title == "" {
				http.Error(w, "service and title are required", http.StatusBadRequest)
				return
			}
			id := s.Add(body.Service, body.Title, body.Message, body.Severity)
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"id":%q}`, id)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/incidents/active", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.Active())
	})

	mux.HandleFunc("/incidents/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 || parts[2] != "resolve" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		id := parts[1]
		if !s.Resolve(id) {
			http.Error(w, "incident not found or already resolved", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
