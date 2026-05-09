package traceid

import (
	"encoding/json"
	"net/http"
	"strings"
)

type traceResponse struct {
	Key string `json:"key"`
	ID  string `json:"id"`
}

// NewHandler returns an HTTP handler that exposes the trace ID store.
// GET  /traceid/{key}  — retrieve stored trace ID
// PUT  /traceid/{key}  — store a trace ID (body: {"id":"..."} or auto-generate)
// DELETE /traceid/{key} — remove stored trace ID
func NewHandler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := strings.TrimPrefix(r.URL.Path, "/traceid/")
		if key == "" {
			http.Error(w, "key required", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodGet:
			id, ok := s.Get(key)
			if !ok {
				http.Error(w, "not found", http.StatusNotFound)
				return
			}
			writeJSON(w, traceResponse{Key: key, ID: id})

		case http.MethodPut:
			var body struct {
				ID string `json:"id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			id := body.ID
			if id == "" {
				var err error
				id, err = Generate()
				if err != nil {
					http.Error(w, "failed to generate ID", http.StatusInternalServerError)
					return
				}
			}
			s.Set(key, id)
			w.WriteHeader(http.StatusCreated)
			writeJSON(w, traceResponse{Key: key, ID: id})

		case http.MethodDelete:
			s.Delete(key)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
