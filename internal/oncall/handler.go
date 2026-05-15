package oncall

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// NewHandler returns an http.Handler that exposes the on-call store.
//
//	GET  /oncall          – list all schedules
//	PUT  /oncall/{svc}    – set schedule for a service
//	DELETE /oncall/{svc}  – remove schedule for a service
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/oncall", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/oncall/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/oncall/")
		if svc == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			sc, ok := s.Get(svc)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, sc)
		case http.MethodPut:
			var body struct {
				Owner string `json:"owner"`
				Phone string `json:"phone"`
				Email string `json:"email"`
				Until string `json:"until"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			until, err := time.Parse(time.RFC3339, body.Until)
			if err != nil {
				http.Error(w, "invalid until: use RFC3339", http.StatusBadRequest)
				return
			}
			if err := s.Set(svc, body.Owner, body.Phone, body.Email, until); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
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

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
