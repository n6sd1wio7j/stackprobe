package audit

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsAllEvents(t *testing.T) {
	s := New(10)
	s.Record(KindCheck, "api", "sys", "up")
	s.Record(KindAlert, "db", "sys", "down")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	NewHandler(s).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var events []Event
	if err := json.NewDecoder(rec.Body).Decode(&events); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(events))
	}
}

func TestHandler_FilterByKind(t *testing.T) {
	s := New(10)
	s.Record(KindCheck, "api", "sys", "check")
	s.Record(KindAlert, "api", "sys", "alert")

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit?kind=check", nil)
	NewHandler(s).ServeHTTP(rec, req)

	var events []Event
	json.NewDecoder(rec.Body).Decode(&events)
	if len(events) != 1 {
		t.Fatalf("expected 1 check event, got %d", len(events))
	}
	if events[0].Kind != KindCheck {
		t.Fatalf("expected kind check, got %q", events[0].Kind)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	s := New(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/audit", nil)
	NewHandler(s).ServeHTTP(rec, req)

	var events []Event
	json.NewDecoder(rec.Body).Decode(&events)
	if len(events) != 0 {
		t.Fatalf("expected empty list, got %d", len(events))
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	s := New(10)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/audit", nil)
	NewHandler(s).ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
