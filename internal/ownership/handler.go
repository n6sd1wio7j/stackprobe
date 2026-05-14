package ownership

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler that exposes the ownership store via REST.
//
//	GET    /ownership          — list all owners
//	GET    /ownership/{svc}    — get owner for service
//	PUT    /ownership/{svc}    — set owner for service
//	DELETE /ownership/{svc}    — remove owner record
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/ownership/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/ownership/")
		switch r.Method {
		case http.MethodGet:
			if service == "" {
				writeJSON(w, http.StatusOK, s.All())
				return
			}
			o, ok := s.Get(service)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, http.StatusOK, o)
		case http.MethodPut:
			var o Owner
			if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			o.Service = service
			if err := s.Set(o); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			writeJSON(w, http.StatusOK, o)
		case http.MethodDelete:
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
