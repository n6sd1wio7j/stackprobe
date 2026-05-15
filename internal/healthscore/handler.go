package healthscore

import (
	"encoding/json"
	"net/http"
	"strings"
)

type handler struct{ store *Store }

// NewHandler returns an HTTP handler for the health-score API.
//
//	GET  /healthscore          – list all scores
//	GET  /healthscore/{service} – get score for one service
func NewHandler(store *Store) http.Handler {
	h := &handler{store: store}
	mux := http.NewServeMux()
	mux.HandleFunc("/healthscore", h.listAll)
	mux.HandleFunc("/healthscore/", h.getOne)
	return mux
}

func (h *handler) listAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, h.store.All())
}

func (h *handler) getOne(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	service := strings.TrimPrefix(r.URL.Path, "/healthscore/")
	if service == "" {
		http.Error(w, "service name required", http.StatusBadRequest)
		return
	}
	sc, ok := h.store.Get(service)
	if !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	writeJSON(w, sc)
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
