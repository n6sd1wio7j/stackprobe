package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newBackend(body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))
}

func TestNew_InvalidBackend(t *testing.T) {
	_, err := New([]Route{{Prefix: "/a", Backend: "://bad url"}})
	if err == nil {
		t.Fatal("expected error for invalid backend URL")
	}
}

func TestServeHTTP_MatchesPrefix(t *testing.T) {
	backend := newBackend("hello")
	defer backend.Close()

	p, err := New([]Route{{Prefix: "/api", Backend: backend.URL}})
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/status", nil)
	p.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello" {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}

func TestServeHTTP_NoMatch(t *testing.T) {
	p, _ := New(nil)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/unknown", nil)
	p.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("expected 502, got %d", rec.Code)
	}
}

func TestAdd_RuntimeRoute(t *testing.T) {
	backend := newBackend("dynamic")
	defer backend.Close()

	p, _ := New(nil)
	if err := p.Add(Route{Prefix: "/dyn", Backend: backend.URL}); err != nil {
		t.Fatalf("Add: %v", err)
	}

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/dyn/test", nil)
	p.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRoutes_ReturnsSnapshot(t *testing.T) {
	p, _ := New([]Route{
		{Prefix: "/a", Backend: "http://localhost:9001"},
		{Prefix: "/b", Backend: "http://localhost:9002"},
	})
	routes := p.Routes()
	if len(routes) != 2 {
		t.Fatalf("expected 2 routes, got %d", len(routes))
	}
}
