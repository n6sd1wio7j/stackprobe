package snapshot

import (
	"encoding/json"
	"net/http"
)

// NewHandler returns an HTTP handler that exposes snapshot data.
//
// GET /snapshots        — returns all latest service snapshots.
// GET /snapshots/{svc}  — returns the snapshot for a single service.
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})

	mux.HandleFunc("/snapshots/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		svc := r.URL.Path[len("/snapshots/"):]
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		rec, ok := s.Get(svc)
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		writeJSON(w, rec)
	})

	return mux
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
