package semver

import (
	"encoding/json"
	"net/http"
	"strings"
)

type handler struct{ store *Store }

// NewHandler returns an HTTP handler for the semver store.
// Routes:
//
//	GET  /semver/          — list all service versions
//	PUT  /semver/{service} — set version (body: {"version":"1.2.3"})
//	DELETE /semver/{service} — remove entry
func NewHandler(s *Store) http.Handler {
	return &handler{store: s}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	service := strings.TrimPrefix(r.URL.Path, "/semver/")

	switch {
	case r.Method == http.MethodGet && service == "":
		h.listAll(w)
	case r.Method == http.MethodPut && service != "":
		h.set(w, r, service)
	case r.Method == http.MethodDelete && service != "":
		h.delete(w, service)
	default:
		http.Error(w, "not found", http.StatusNotFound)
	}
}

func (h *handler) listAll(w http.ResponseWriter) {
	all := h.store.All()
	out := make(map[string]string, len(all))
	for k, v := range all {
		out[k] = v.Raw
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
}

func (h *handler) set(w http.ResponseWriter, r *http.Request, service string) {
	var body struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Version == "" {
		http.Error(w, "invalid body: expected {\"version\":\"x.y.z\"}", http.StatusBadRequest)
		return
	}
	if err := h.store.Set(service, body.Version); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) delete(w http.ResponseWriter, service string) {
	h.store.Delete(service)
	w.WriteHeader(http.StatusNoContent)
}
