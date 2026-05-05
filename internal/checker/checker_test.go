package checker_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stackprobe/stackprobe/internal/checker"
)

func newTestServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
	}))
}

func TestCheck_Up(t *testing.T) {
	srv := newTestServer(http.StatusOK)
	defer srv.Close()

	c := checker.New(5 * time.Second)
	ep := checker.Endpoint{Name: "test-ok", URL: srv.URL}

	res := c.Check(context.Background(), ep)

	if res.Status != checker.StatusUp {
		t.Errorf("expected status %q, got %q", checker.StatusUp, res.Status)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status code 200, got %d", res.StatusCode)
	}
	if res.Error != "" {
		t.Errorf("expected no error, got %q", res.Error)
	}
}

func TestCheck_Down_ServerError(t *testing.T) {
	srv := newTestServer(http.StatusInternalServerError)
	defer srv.Close()

	c := checker.New(5 * time.Second)
	ep := checker.Endpoint{Name: "test-500", URL: srv.URL}

	res := c.Check(context.Background(), ep)

	if res.Status != checker.StatusDown {
		t.Errorf("expected status %q, got %q", checker.StatusDown, res.Status)
	}
}

func TestCheck_Down_InvalidURL(t *testing.T) {
	c := checker.New(2 * time.Second)
	ep := checker.Endpoint{Name: "bad-url", URL: "http://127.0.0.1:0/nonexistent"}

	res := c.Check(context.Background(), ep)

	if res.Status != checker.StatusDown {
		t.Errorf("expected status %q, got %q", checker.StatusDown, res.Status)
	}
	if res.Error == "" {
		t.Error("expected an error message but got none")
	}
}

func TestCheckAll_ConcurrentResults(t *testing.T) {
	srv1 := newTestServer(http.StatusOK)
	defer srv1.Close()
	srv2 := newTestServer(http.StatusServiceUnavailable)
	defer srv2.Close()

	c := checker.New(5 * time.Second)
	endpoints := []checker.Endpoint{
		{Name: "svc-a", URL: srv1.URL},
		{Name: "svc-b", URL: srv2.URL},
	}

	results := c.CheckAll(context.Background(), endpoints)

	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Status != checker.StatusUp {
		t.Errorf("svc-a: expected %q, got %q", checker.StatusUp, results[0].Status)
	}
	if results[1].Status != checker.StatusDown {
		t.Errorf("svc-b: expected %q, got %q", checker.StatusDown, results[1].Status)
	}
}
