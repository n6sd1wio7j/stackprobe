package endpoint

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler for the endpoint registry.
// Routes:
//
//	GET    /endpoints          – list all
//	PUT    /endpoints/{svc}    – register / update
//	GET    /endpoints/{svc}    – get one
//	DELETE /endpoints/{svc}    – remove
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/endpoints", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All(), http.StatusOK)
	})
	mux.HandleFunc("/endpoints/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/endpoints/")
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			m, err := s.Get(svc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, m, http.StatusOK)
		case http.MethodPut:
			var body struct {
				URL         string `json:"url"`
				Description string `json:"description"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Register(svc, body.URL, body.Description); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			m, _ := s.Get(svc)
			writeJSON(w, m, http.StatusOK)
		case http.MethodDelete:
			if err := s.Delete(svc); err != nil {
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

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
