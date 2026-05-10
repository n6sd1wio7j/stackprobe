package labels

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler that exposes label CRUD operations.
//
// Routes:
//
//	GET    /labels/{service}           — retrieve all labels for a service
//	PUT    /labels/{service}/{key}      — set a label value (JSON body: {"value":"..."})
//	DELETE /labels/{service}/{key}      — remove a label
//	GET    /labels                      — list all services and their labels
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/labels", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})

	mux.HandleFunc("/labels/", func(w http.ResponseWriter, r *http.Request) {
		// strip leading "/labels/"
		parts := strings.SplitN(strings.TrimPrefix(r.URL.Path, "/labels/"), "/", 2)
		service := parts[0]
		if service == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			// /labels/{service}
			if r.Method != http.MethodGet {
				http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
				return
			}
			labels, err := s.Get(service)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, labels)
			return
		}

		// /labels/{service}/{key}
		key := parts[1]
		switch r.Method {
		case http.MethodPut:
			var body struct {
				Value string `json:"value"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			s.Set(service, key, body.Value)
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			s.Delete(service, key)
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
