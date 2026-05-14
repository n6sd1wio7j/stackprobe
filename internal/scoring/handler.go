package scoring

import (
	"encoding/json"
	"net/http"
	"strings"
)

// NewHandler returns an HTTP handler for the scoring store.
// Routes:
//   GET  /scoring          – list all scores
//   GET  /scoring/{svc}    – get score for a service
//   PUT  /scoring/{svc}    – set score (body: {"uptime_pct": 99.5, "avg_latency_ms": 45})
//   DELETE /scoring/{svc}  – remove score
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/scoring", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, s.All(), http.StatusOK)
	})

	mux.HandleFunc("/scoring/", func(w http.ResponseWriter, r *http.Request) {
		svc := strings.TrimPrefix(r.URL.Path, "/scoring/")
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
			writeJSON(w, sc, http.StatusOK)

		case http.MethodPut:
			var body struct {
				UptimePct    float64 `json:"uptime_pct"`
				AvgLatencyMs float64 `json:"avg_latency_ms"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			if err := s.Set(svc, body.UptimePct, body.AvgLatencyMs); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			sc, _ := s.Get(svc)
			writeJSON(w, sc, http.StatusOK)

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
