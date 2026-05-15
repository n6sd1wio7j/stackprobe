package catalog

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the catalog store.
// Routes:
//
//	GET    /catalog          -> list all entries
//	PUT    /catalog/{svc}    -> register/update entry
//	GET    /catalog/{svc}    -> get single entry
//	DELETE /catalog/{svc}    -> remove entry
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/catalog", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/catalog/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/catalog/")
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			e, ok := s.Get(svc)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, e)
		case http.MethodPut:
			var body struct {
				Description string `json:"description"`
				Owner       string `json:"owner"`
				URL         string `json:"url"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Register(svc, body.Description, body.Owner, body.URL); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			e, _ := s.Get(svc)
			w.WriteHeader(http.StatusOK)
			writeJSON(w, e)
		case http.MethodDelete:
			if err := s.Delete(svc); err != nil {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
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
