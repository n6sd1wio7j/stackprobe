package routing

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler that exposes routing rule management.
//
//	GET  /routing        → list all rules
//	POST /routing        → add a rule   (body: {"prefix":"/x","target":"svc","strip_prefix":false})
//	DELETE /routing?prefix=/x → remove rule by prefix
func NewHandler(r *Router) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodGet:
			rules := r.All()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(rules)

		case http.MethodPost:
			var rule Rule
			if err := json.NewDecoder(req.Body).Decode(&rule); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}
			if err := r.Add(rule); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)

		case http.MethodDelete:
			prefix := strings.TrimSpace(req.URL.Query().Get("prefix"))
			if prefix == "" {
				http.Error(w, "prefix query parameter required", http.StatusBadRequest)
				return
			}
			r.Remove(prefix)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
