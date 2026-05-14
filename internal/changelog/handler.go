package changelog

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the changelog store.
//
// Routes:
//
//	GET  /changelog            – all entries (optional ?kind= filter)
//	GET  /changelog/{service}  – entries for a specific service
//	POST /changelog/{service}  – add an entry
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/changelog", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		kind := EntryKind(r.URL.Query().Get("kind"))
		writeJSON(w, s.Filter(kind))
	})
	mux.HandleFunc("/changelog/", func(w http.ResponseWriter, r *http.Request) {
		service := strings.TrimPrefix(r.URL.Path, "/changelog/")
		if service == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, s.Get(service))
		case http.MethodPost:
			var body struct {
				Kind    EntryKind `json:"kind"`
				Message string    `json:"message"`
				Author  string    `json:"author"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			e, err := s.Add(service, body.Kind, body.Message, body.Author)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
			writeJSON(w, e)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
