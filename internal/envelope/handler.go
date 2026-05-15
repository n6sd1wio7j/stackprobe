package envelope

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler for envelope management.
// Routes:
//   GET  /envelope/          - list all envelopes
//   PUT  /envelope/{service} - wrap / update an envelope
//   GET  /envelope/{service} - get a single envelope
//   DELETE /envelope/{service} - remove an envelope
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/envelope/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/envelope/")
		switch {
		case service == "" && r.Method == http.MethodGet:
			writeJSON(w, s.All(), http.StatusOK)
		case service != "" && r.Method == http.MethodGet:
			e, ok := s.Get(service)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, e, http.StatusOK)
		case service != "" && r.Method == http.MethodPut:
			var body struct {
				Status  string `json:"status"`
				Message string `json:"message"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Wrap(service, body.Status, body.Message); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			e, _ := s.Get(service)
			writeJSON(w, e, http.StatusOK)
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

func writeJSON(w http.ResponseWriter, v any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
