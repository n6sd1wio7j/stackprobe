package metrics_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stackprobe/internal/metrics"
)

func TestHandler_ReturnsJSON(t *testing.T) {
	store := metrics.New()
	store.Record("web", true, 80*time.Millisecond)
	store.Record("web", false, 120*time.Millisecond)

	h := metrics.NewHandler(store)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(result))
	}
	if result[0]["service"] != "web" {
		t.Errorf("expected service=web, got %v", result[0]["service"])
	}
	if result[0]["total_checks"].(float64) != 2 {
		t.Errorf("expected total_checks=2, got %v", result[0]["total_checks"])
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	store := metrics.New()
	h := metrics.NewHandler(store)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/metrics", nil)
	h(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rr.Code)
	}
}

func TestHandler_EmptyStore(t *testing.T) {
	store := metrics.New()
	h := metrics.NewHandler(store)
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	h(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d entries", len(result))
	}
}
