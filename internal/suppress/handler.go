package suppress

import (
	"encoding/json"
	"net/http"
	"time"
)

type handler struct{ store *Store }

// NewHandler returns an HTTP handler for the suppression API.
//
//	GET  /suppress          – list all rules
//	PUT  /suppress/{service} – add/replace a rule  (body: {"reason":"…","until":"RFC3339"})
//	DELETE /suppress/{service} – remove a rule
func NewHandler(s *Store) http.Handler {
	mux := http.NewServeMux()
	h := &handler{store: s}
	mux.HandleFunc("/suppress", h.listAll)
	mux.HandleFunc("/suppress/", h.dispatch)
	return mux
}

func (h *handler) listAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, h.store.All())
}

func (h *handler) dispatch(w http.ResponseWriter, r *http.Request) {
	service := r.URL.Path[len("/suppress/"):]
	if service == "" {
		http.Error(w, "service required", http.StatusBadRequest)
		return
	}
	switch r.Method {
	case http.MethodPut:
		var body struct {
			Reason string `json:"reason"`
			Until  string `json:"until"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}
		until, err := time.Parse(time.RFC3339, body.Until)
		if err != nil {
			http.Error(w, "invalid until: use RFC3339", http.StatusBadRequest)
			return
		}
		if err := h.store.Suppress(service, body.Reason, until); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	case http.MethodDelete:
		if err := h.store.Unsuppress(service); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
