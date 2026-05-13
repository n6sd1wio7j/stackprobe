package weight

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// NewHandler returns an HTTP handler for the weight store.
//
//	GET  /weights          — list all explicit weights
//	PUT  /weights/{svc}    — set weight (JSON body: {"weight": N})
//	DELETE /weights/{svc}  — remove explicit weight
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/weights", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.All())
	})
	mux.HandleFunc("/weights/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/weights/")
		if svc == "" {
			http.Error(w, "service required", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]int{"weight": s.Get(svc)})
		case http.MethodPut:
			var body struct {
				Weight int `json:"weight"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(svc, body.Weight); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"service": svc, "weight": strconv.Itoa(body.Weight)})
		case http.MethodDelete:
			s.Delete(svc)
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
	return mux
}
