package heartbeat

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newHandler(t *testing.T) (*Store, http.Handler) {
	t.Helper()
	s := New(time.Minute)
	return s, NewHandler(s)
}

func TestHandler_Ping_RecordsBeat(t *testing.T) {
	store, h := newHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/heartbeat/api", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rw.Code)
	}
	_, ok := store.Get("api")
	if !ok {
		t.Error("expected beat to be recorded")
	}
}

func TestHandler_Ping_MissingService(t *testing.T) {
	_, h := newHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/heartbeat/", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rw.Code)
	}
}

func TestHandler_GetAll(t *testing.T) {
	store, h := newHandler(t)
	store.Ping("svc1")
	store.Ping("svc2")
	req := httptest.NewRequest(http.MethodGet, "/heartbeat", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var beats []Beat
	if err := json.NewDecoder(rw.Body).Decode(&beats); err != nil {
		t.Fatal(err)
	}
	if len(beats) != 2 {
		t.Errorf("expected 2 beats, got %d", len(beats))
	}
}

func TestHandler_GetStale(t *testing.T) {
	store, h := newHandler(t)
	store.Register("slow", time.Nanosecond)
	store.Ping("slow")
	time.Sleep(time.Millisecond)

	req := httptest.NewRequest(http.MethodGet, "/heartbeat/stale", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}
	var beats []Beat
	if err := json.NewDecoder(rw.Body).Decode(&beats); err != nil {
		t.Fatal(err)
	}
	if len(beats) == 0 {
		t.Error("expected at least one stale beat")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	_, h := newHandler(t)
	req := httptest.NewRequest(http.MethodDelete, "/heartbeat", nil)
	rw := httptest.NewRecorder()
	h.ServeHTTP(rw, req)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rw.Code)
	}
}
