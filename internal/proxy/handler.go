package proxy

import (
	"encoding/json"
	"net/http"
)

type routeRequest struct {
	Prefix  string `json:"prefix"`
	Backend string `json:"backend"`
}

// NewHandler returns an HTTP handler for managing proxy routes via REST.
// GET  /proxy/routes        — list all routes
// POST /proxy/routes        — add a new route
func NewHandler(p *Proxy) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/proxy/routes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			routes := p.Routes()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(routes)

		case http.MethodPost:
			var req routeRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}
			if req.Prefix == "" || req.Backend == "" {
				http.Error(w, "prefix and backend are required", http.StatusBadRequest)
				return
			}
			if err := p.Add(Route{Prefix: req.Prefix, Backend: req.Backend}); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
