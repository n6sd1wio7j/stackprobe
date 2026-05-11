package region

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the region store.
// Routes:
//   GET  /regions          -> all services grouped by region
//   GET  /regions/{svc}    -> region for a specific service
//   PUT  /regions/{svc}    -> assign region to service  {"region":"..."}
//   DELETE /regions/{svc}  -> remove region assignment
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/regions", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(s.ByRegion())
	})

	mux.HandleFunc("/regions/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/regions/")
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			reg, ok := s.Get(svc)
			if !ok {
				http.Error(w, "service not found", http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"service": svc, "region": reg})

		case http.MethodPut:
			var body struct {
				Region string `json:"region"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Region == "" {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			if err := s.Set(svc, body.Region); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			s.Delete(svc)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
