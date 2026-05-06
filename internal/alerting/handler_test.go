package alerting

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_ReturnsActiveAlerts(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", false)
	m.Evaluate("svc-b", false)

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rec := httptest.NewRecorder()
	NewHandler(m).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var alerts []alertResponse
	if err := json.NewDecoder(rec.Body).Decode(&alerts); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(alerts) != 2 {
		t.Errorf("expected 2 alerts, got %d", len(alerts))
	}
}

func TestHandler_EmptyWhenNoAlerts(t *testing.T) {
	m := New()

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rec := httptest.NewRecorder()
	NewHandler(m).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var alerts []alertResponse
	if err := json.NewDecoder(rec.Body).Decode(&alerts); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(alerts) != 0 {
		t.Errorf("expected 0 alerts, got %d", len(alerts))
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	m := New()

	req := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	rec := httptest.NewRecorder()
	NewHandler(m).ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
