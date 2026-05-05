package notifier_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stackprobe/internal/notifier"
)

func newWebhookServer(t *testing.T, statusCode int, received *notifier.AlertPayload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, received)
		w.WriteHeader(statusCode)
	}))
}

func TestNotify_Success(t *testing.T) {
	var got notifier.AlertPayload
	srv := newWebhookServer(t, http.StatusOK, &got)
	defer srv.Close()

	n := notifier.New(srv.URL, 5*time.Second)
	if err := n.Notify("api", "down", "connection refused"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got.Service != "api" {
		t.Errorf("service: want %q, got %q", "api", got.Service)
	}
	if got.Status != "down" {
		t.Errorf("status: want %q, got %q", "down", got.Status)
	}
	if got.Message != "connection refused" {
		t.Errorf("message: want %q, got %q", "connection refused", got.Message)
	}
	if got.Timestamp.IsZero() {
		t.Error("expected non-zero timestamp")
	}
}

func TestNotify_Non2xxReturnsError(t *testing.T) {
	var got notifier.AlertPayload
	srv := newWebhookServer(t, http.StatusInternalServerError, &got)
	defer srv.Close()

	n := notifier.New(srv.URL, 5*time.Second)
	if err := n.Notify("db", "up", ""); err == nil {
		t.Fatal("expected error for non-2xx response, got nil")
	}
}

func TestNotify_InvalidURL(t *testing.T) {
	n := notifier.New("http://127.0.0.1:0", 1*time.Second)
	if err := n.Notify("svc", "down", "timeout"); err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}
