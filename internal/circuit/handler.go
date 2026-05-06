package circuit

import (
	"encoding/json"
	"net/http"
)

type statusResponse struct {
	Service string `json:"service"`
	State   string `json:"state"`
}

// NewHandler returns an HTTP handler that exposes circuit breaker states.
// GET /circuit  — returns JSON array of all tracked services and their states.
func NewHandler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		s.mu.Lock()
		services := make([]statusResponse, 0, len(s.breakers))
		for name, b := range s.breakers {
			b.mu.Lock()
			services = append(services, statusResponse{
				Service: name,
				State:   b.state.String(),
			})
			b.mu.Unlock()
		}
		s.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(services); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
	})
}
