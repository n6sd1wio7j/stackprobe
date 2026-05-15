package badge

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) (http.Handler, *Store) {
	t.Helper()
	s := New()
	return NewHandler(s), s
}

func TestHandler_GetAll_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/badge/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Put_CreatesBadge(t *testing.T) {
	h, s := newHandler(t)
	body, _ := json.Marshal(map[string]string{
		"label": "status", "status": "up", "color": "green", "style": "flat",
	})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/badge/api", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if _, err := s.Get("api"); err != nil {
		t.Fatalf("badge not stored: %v", err)
	}
}

func TestHandler_Put_InvalidJSON(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/badge/api", bytes.NewBufferString("not-json")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Get_KnownService(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("api", "status", "up", "green", StyleFlat)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/badge/api", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var b Badge
	if err := json.NewDecoder(rec.Body).Decode(&b); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if b.Status != "up" {
		t.Errorf("expected status up, got %q", b.Status)
	}
}

func TestHandler_Get_UnknownService(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/badge/ghost", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_Delete_RemovesBadge(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("api", "status", "up", "green", StyleFlat)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/badge/api", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/badge/", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
