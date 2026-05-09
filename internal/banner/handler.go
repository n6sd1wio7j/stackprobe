package banner

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// NewHandler returns an http.Handler for banner management.
// GET  /banners        — list active banners
// POST /banners        — create a banner
// DELETE /banners/{id} — remove a banner
func NewHandler(s *Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.URL.Path == "/banners" || r.URL.Path == "/banners/" {
			switch r.Method {
			case http.MethodGet:
				listBanners(w, s)
			case http.MethodPost:
				createBanner(w, r, s)
			default:
				http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
			}
			return
		}

		// DELETE /banners/{id}
		id := strings.TrimPrefix(r.URL.Path, "/banners/")
		if r.Method == http.MethodDelete && id != "" {
			if s.Remove(id) {
				w.WriteHeader(http.StatusNoContent)
			} else {
				http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
			}
			return
		}

		http.Error(w, `{"error":"not found"}`, http.StatusNotFound)
	})
}

func listBanners(w http.ResponseWriter, s *Store) {
	active := s.Active()
	if active == nil {
		active = []*Banner{}
	}
	json.NewEncoder(w).Encode(active)
}

func createBanner(w http.ResponseWriter, r *http.Request, s *Store) {
	var req struct {
		Message   string `json:"message"`
		Level     string `json:"level"`
		ExpiresIn string `json:"expires_in"` // e.g. "2h"
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}
	if req.Message == "" {
		http.Error(w, `{"error":"message required"}`, http.StatusBadRequest)
		return
	}
	if req.Level == "" {
		req.Level = "info"
	}

	var expiresAt time.Time
	if req.ExpiresIn != "" {
		d, err := time.ParseDuration(req.ExpiresIn)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"invalid expires_in: %s"}`, err), http.StatusBadRequest)
			return
		}
		expiresAt = time.Now().Add(d)
	}

	id := s.Add(req.Message, req.Level, expiresAt)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}
