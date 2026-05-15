package badge

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler for the badge API.
//
// Routes:
//   GET    /badge/          -> list all badges
//   PUT    /badge/{service} -> create or update a badge
//   GET    /badge/{service} -> get a single badge
//   DELETE /badge/{service} -> remove a badge
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/badge/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/badge/")
		switch {
		case service == "" && r.Method == http.MethodGet:
			writeJSON(w, http.StatusOK, s.All())
		case service != "" && r.Method == http.MethodGet:
			b, err := s.Get(service)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, b)
		case service != "" && r.Method == http.MethodPut:
			var req struct {
				Label  string `json:"label"`
				Status string `json:"status"`
				Color  string `json:"color"`
				Style  Style  `json:"style"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(service, req.Label, req.Status, req.Color, req.Style); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			b, _ := s.Get(service)
			writeJSON(w, http.StatusOK, b)
		case service != "" && r.Method == http.MethodDelete:
			if err := s.Delete(service); err != nil {
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

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
