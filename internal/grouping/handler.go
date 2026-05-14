package grouping

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an http.Handler that exposes grouping over HTTP.
//
//	GET    /groups          – list all groups
//	PUT    /groups/{name}   – create/replace group
//	GET    /groups/{name}   – get a single group
//	DELETE /groups/{name}   – delete a group
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All())
	})
	mux.HandleFunc("/groups/", func(w http.ResponseWriter, r *http.Request) {
		name := strings.TrimPrefix(r.URL.Path, "/groups/")
		if name == "" {
			http.Error(w, "group name required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			svcs, ok := s.Get(name)
			if !ok {
				http.Error(w, "group not found", http.StatusNotFound)
				return
			}
			writeJSON(w, Group{Name: name, Services: svcs})
		case http.MethodPut:
			var body struct {
				Services []string `json:"services"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(name, body.Services); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		case http.MethodDelete:
			if err := s.Delete(name); err != nil {
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
