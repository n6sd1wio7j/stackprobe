package environ

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler that exposes the environ Store via REST.
//
//	GET  /environ           → all mappings
//	GET  /environ/{service} → single service
//	PUT  /environ/{service} → set environment  (body: {"env":"production"})
//	DELETE /environ/{service} → remove mapping
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/environ", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/environ/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/environ/")
		if service == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			env, ok := s.Get(service)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, map[string]string{"service": service, "env": env})
		case http.MethodPut:
			var body struct {
				Env string `json:"env"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Env == "" {
				http.Error(w, "invalid body: env required", http.StatusBadRequest)
				return
			}
			if err := s.Set(service, body.Env); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			s.Delete(service)
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
