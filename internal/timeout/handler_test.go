package timeout_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/your-org/stackprobe/internal/timeout"
)

func newHandler(t *testing.T) (http.Handler, *timeout.Store) {
	t.Helper()
	s := timeout.New(5 * time.Second)
	return timeout.NewHandler(s), s
}

func TestHandler_GetAll(t *testing.T) {
	h, s := newHandler(t)
	s.Set("svc-a", 2*time.Second)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out["svc-a"] != "2s" {
		t.Errorf("expected 2s for svc-a, got %q", out["svc-a"])
	}
}

func TestHandler_Put_SetsOverride(t *testing.T) {
	h, s := newHandler(t)
	body := bytes.NewBufferString(`{"timeout":"3s"}`)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/svc-b", body))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if s.Get("svc-b") != 3*time.Second {
		t.Errorf("expected 3s for svc-b, got %v", s.Get("svc-b"))
	}
}

func TestHandler_Put_InvalidDuration(t *testing.T) {
	h, _ := newHandler(t)
	body := bytes.NewBufferString(`{"timeout":"not-a-duration"}`)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/svc-c", body))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Delete_RemovesOverride(t *testing.T) {
	h, s := newHandler(t)
	s.Set("svc-d", 1*time.Second)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/svc-d", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if s.Get("svc-d") != 5*time.Second {
		t.Errorf("expected default after delete, got %v", s.Get("svc-d"))
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
