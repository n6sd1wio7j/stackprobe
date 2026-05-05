package history

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Handler exposes the history store over HTTP.
// GET /history          → returns all service histories
// GET /history/{service} → returns history for a single service
type Handler struct {
	store *Store
}

// NewHandler wraps a Store in an HTTP handler.
func NewHandler(store *Store) *Handler {
	return &Handler{store: store}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Trim the base prefix so callers can mount at any path.
	path := strings.TrimPrefix(r.URL.Path, "/history")
	path = strings.TrimPrefix(path, "/")

	w.Header().Set("Content-Type", "application/json")

	if path == "" {
		// Return all service histories.
		all := h.store.All()
		if err := json.NewEncoder(w).Encode(all); err != nil {
			http.Error(w, "encoding error", http.StatusInternalServerError)
		}
		return
	}

	// Return history for the named service.
	recs := h.store.Get(path)
	if err := json.NewEncoder(w).Encode(recs); err != nil {
		http.Error(w, "encoding error", http.StatusInternalServerError)
	}
}
