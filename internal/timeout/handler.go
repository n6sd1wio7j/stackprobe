package timeout

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type handler struct {
	store *Store
}

// NewHandler returns an http.Handler that exposes timeout overrides via REST.
//
//	GET  /timeouts          – list all overrides
//	PUT  /timeouts/{svc}    – set override (body: {"timeout":"2s"})
//	DELETE /timeouts/{svc}  – remove override
func NewHandler(s *Store) http.Handler {
	return &handler{store: s}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	svc := strings.TrimPrefix(r.URL.Path, "/")

	switch r.Method {
	case http.MethodGet:
		all := h.store.All()
		out := make(map[string]string, len(all))
		for k, v := range all {
			out[k] = v.String()
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out) //nolint:errcheck

	case http.MethodPut:
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		var body struct {
			Timeout string `json:"timeout"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		d, err := time.ParseDuration(body.Timeout)
		if err != nil || d <= 0 {
			http.Error(w, "invalid timeout duration", http.StatusBadRequest)
			return
		}
		h.store.Set(svc, d)
		w.WriteHeader(http.StatusNoContent)

	case http.MethodDelete:
		if svc == "" {
			http.Error(w, "service name required", http.StatusBadRequest)
			return
		}
		h.store.Delete(svc)
		w.WriteHeader(http.StatusNoContent)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
