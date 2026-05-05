package dashboard_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stackprobe/internal/aggregator"
	"github.com/stackprobe/internal/dashboard"
)

// stubCollector implements dashboard.Collector for testing.
type stubCollector struct {
	results []aggregator.ServiceStatus
}

func (s *stubCollector) Collect() []aggregator.ServiceStatus {
	return s.results
}

func TestHandler_AllHealthy(t *testing.T) {
	collector := &stubCollector{
		results: []aggregator.ServiceStatus{
			{Name: "api", URL: "http://api", Up: true, CheckedAt: time.Now()},
			{Name: "db", URL: "http://db", Up: true, CheckedAt: time.Now()},
		},
	}

	h := dashboard.New(collector)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var status dashboard.Status
	if err := json.NewDecoder(rec.Body).Decode(&status); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if status.Healthy != 2 || status.Unhealthy != 0 {
		t.Errorf("expected 2 healthy / 0 unhealthy, got %d / %d", status.Healthy, status.Unhealthy)
	}
	if len(status.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(status.Services))
	}
}

func TestHandler_PartialFailure(t *testing.T) {
	collector := &stubCollector{
		results: []aggregator.ServiceStatus{
			{Name: "api", URL: "http://api", Up: true, CheckedAt: time.Now()},
			{Name: "cache", URL: "http://cache", Up: false, CheckedAt: time.Now()},
		},
	}

	h := dashboard.New(collector)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	h.ServeHTTP(rec, req)

	var status dashboard.Status
	_ = json.NewDecoder(rec.Body).Decode(&status)

	if status.Healthy != 1 || status.Unhealthy != 1 {
		t.Errorf("expected 1 healthy / 1 unhealthy, got %d / %d", status.Healthy, status.Unhealthy)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := dashboard.New(&stubCollector{})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}
