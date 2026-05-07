package proxy

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) (http.Handler, *Proxy) {
	t.Helper()
	p, _ := New(nil)
	return NewHandler(p), p
}

func TestHandler_GetRoutes_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/proxy/routes", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var routes []Route
	_ = json.NewDecoder(rec.Body).Decode(&routes)
	if len(routes) != 0 {
		t.Fatalf("expected empty list")
	}
}

func TestHandler_PostRoute_Valid(t *testing.T) {
	h, p := newHandler(t)
	body, _ := json.Marshal(routeRequest{Prefix: "/svc", Backend: "http://localhost:8080"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/proxy/routes", bytes.NewReader(body))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
	if len(p.Routes()) != 1 {
		t.Fatal("route not added")
	}
}

func TestHandler_PostRoute_MissingFields(t *testing.T) {
	h, _ := newHandler(t)
	body, _ := json.Marshal(routeRequest{Prefix: "/only"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/proxy/routes", bytes.NewReader(body))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_PostRoute_InvalidBackend(t *testing.T) {
	h, _ := newHandler(t)
	body, _ := json.Marshal(routeRequest{Prefix: "/x", Backend: "://bad"})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/proxy/routes", bytes.NewReader(body))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/proxy/routes", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
