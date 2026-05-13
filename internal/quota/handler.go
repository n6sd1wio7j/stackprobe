package quota

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the quota store.
//
// Routes:
//
//	GET  /quota/          — list all quota entries
//	PUT  /quota/{service} — set quota limit (body: {"limit": N})
//	GET  /quota/{service} — get quota entry for service
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/quota/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/quota/")
		switch {
		case service == "" && r.Method == http.MethodGet:
			writeJSON(w, s.All())
		case service != "" && r.Method == http.MethodPut:
			var body struct {
				Limit int `json:"limit"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(service, body.Limit); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case service != "" && r.Method == http.MethodGet:
			e, err := s.Get(service)
			if err == ErrUnknownService {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, e)
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
