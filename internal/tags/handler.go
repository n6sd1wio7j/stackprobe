package tags

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler that exposes tag query and management.
// GET  /tags?filter=tag1,tag2  — list services matching all given tags
// GET  /tags                   — return all service->tag mappings
// POST /tags/{service}         — set tags for a service (JSON body: ["tag1","tag2"])
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if f := r.URL.Query().Get("filter"); f != "" {
			required := strings.Split(f, ",")
			matches := s.Filter(required...)
			if matches == nil {
				matches = []string{}
			}
			json.NewEncoder(w).Encode(matches)
			return
		}
		json.NewEncoder(w).Encode(s.All())
	})

	mux.HandleFunc("/tags/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		service := strings.TrimPrefix(r.URL.Path, "/tags/")
		if service == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		var incoming []string
		if err := json.NewDecoder(r.Body).Decode(&incoming); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		s.Set(service, incoming)
		w.WriteHeader(http.StatusNoContent)
	})

	return mux
}
