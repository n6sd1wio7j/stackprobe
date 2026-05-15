package checkpoint

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the checkpoint store.
//
//	GET    /checkpoint          → list all
//	PUT    /checkpoint/{svc}    → save checkpoint
//	GET    /checkpoint/{svc}    → get checkpoint
//	DELETE /checkpoint/{svc}    → delete checkpoint
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/checkpoint", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/checkpoint/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/checkpoint/")
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			rec, err := s.Get(svc)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			writeJSON(w, rec)
		case http.MethodPut:
			var body struct {
				Status string `json:"status"`
				Note   string `json:"note"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Save(svc, body.Status, body.Note); err != nil {
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
