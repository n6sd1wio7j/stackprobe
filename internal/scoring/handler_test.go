package scoring

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
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/scoring", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
}

func TestHandler_Put_SetsScore(t *testing.T) {
	h, s := newHandler(t)
	body := `{"uptime_pct":98.5,"avg_latency_ms":60}`
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/scoring/api", bytes.NewBufferString(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	sc, ok := s.Get("api")
	if !ok {
		t.Fatal("expected score to be stored")
	}
	if sc.Grade == "" {
		t.Error("expected non-empty grade")
	}
}

func TestHandler_Put_InvalidJSON(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/scoring/api", bytes.NewBufferString("bad")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("got %d, want 400", rec.Code)
	}
}

func TestHandler_Get_KnownService(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("db", 99, 20)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/scoring/db", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", rec.Code)
	}
	var sc Score
	if err := json.NewDecoder(rec.Body).Decode(&sc); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if sc.Service != "db" {
		t.Errorf("got service %q, want db", sc.Service)
	}
}

func TestHandler_Get_UnknownService(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/scoring/ghost", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("got %d, want 404", rec.Code)
	}
}

func TestHandler_Delete_RemovesScore(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("cache", 95, 80)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/scoring/cache", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("got %d, want 204", rec.Code)
	}
	if _, ok := s.Get("cache"); ok {
		t.Fatal("expected score to be removed")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/scoring", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got %d, want 405", rec.Code)
	}
}
