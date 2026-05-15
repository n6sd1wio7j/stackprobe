package endpoint

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
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/endpoints", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []Meta
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d items", len(out))
	}
}

func TestHandler_Put_RegistersEndpoint(t *testing.T) {
	h, s := newHandler(t)
	body, _ := json.Marshal(map[string]string{"url": "http://svc/health", "description": "test"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/endpoints/svc", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if _, err := s.Get("svc"); err != nil {
		t.Fatal("expected endpoint to be stored")
	}
}

func TestHandler_Put_InvalidJSON(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/endpoints/svc", bytes.NewBufferString("bad")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Get_KnownService(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Register("svc", "http://svc", "desc")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/endpoints/svc", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var m Meta
	_ = json.NewDecoder(rec.Body).Decode(&m)
	if m.URL != "http://svc" {
		t.Errorf("unexpected url %q", m.URL)
	}
}

func TestHandler_Get_UnknownService(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/endpoints/ghost", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_Delete_RemovesEndpoint(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Register("svc", "http://svc", "")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/endpoints/svc", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/endpoints", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
