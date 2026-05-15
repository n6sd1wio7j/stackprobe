package capacity

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler for the capacity API.
//
// Routes:
//   GET    /capacity          → list all records
//   PUT    /capacity/{svc}    → set capacity
//   GET    /capacity/{svc}    → get one record
//   DELETE /capacity/{svc}    → remove record
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/capacity", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All(), http.StatusOK)
	})
	mux.HandleFunc("/capacity/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/capacity/")
		if svc == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			rec, ok := s.Get(svc)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, rec, http.StatusOK)
		case http.MethodPut:
			var body struct {
				Limit   int `json:"limit"`
				Current int `json:"current"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(svc, body.Limit, body.Current); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			rec, _ := s.Get(svc)
			writeJSON(w, rec, http.StatusOK)
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

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
