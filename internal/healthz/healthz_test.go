package healthz_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackprobe/internal/healthz"
)

func TestLiveness_ReturnsOK(t *testing.T) {
	h := healthz.New("v1.2.3")

	req := httptest.NewRequest(http.MethodGet, "/healthz/live", nil)
	rec := httptest.NewRecorder()

	h.LivenessHandler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp healthz.Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status ok, got %q", resp.Status)
	}
	if resp.Version != "v1.2.3" {
		t.Errorf("expected version v1.2.3, got %q", resp.Version)
	}
	if resp.GoVersion == "" {
		t.Error("expected non-empty go_version")
	}
}

func TestLiveness_MethodNotAllowed(t *testing.T) {
	h := healthz.New("")

	req := httptest.NewRequest(http.MethodPost, "/healthz/live", nil)
	rec := httptest.NewRecorder()

	h.LivenessHandler(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestReadiness_WhenReady(t *testing.T) {
	h := healthz.New("v0.1.0")
	handler := h.ReadinessHandler(func() bool { return true })

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp healthz.Response
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Status != "ok" {
		t.Errorf("expected ok, got %q", resp.Status)
	}
}

func TestReadiness_WhenNotReady(t *testing.T) {
	h := healthz.New("v0.1.0")
	handler := h.ReadinessHandler(func() bool { return false })

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected 503, got %d", rec.Code)
	}

	var resp healthz.Response
	_ = json.NewDecoder(rec.Body).Decode(&resp)
	if resp.Status != "unavailable" {
		t.Errorf("expected unavailable, got %q", resp.Status)
	}
}

func TestReadiness_NilReadyFunc(t *testing.T) {
	h := healthz.New("")
	handler := h.ReadinessHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/healthz/ready", nil)
	rec := httptest.NewRecorder()

	handler(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 with nil ready func, got %d", rec.Code)
	}
}
