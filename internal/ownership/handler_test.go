package ownership

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler() http.Handler {
	return NewHandler(New())
}

func TestHandler_GetAll_Empty(t *testing.T) {
	h := newHandler()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ownership/", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var out []Owner
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d items", len(out))
	}
}

func TestHandler_Put_SetsOwner(t *testing.T) {
	h := newHandler()
	body, _ := json.Marshal(map[string]string{"team": "sre", "email": "sre@example.com", "slack": "#sre"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/ownership/payments", bytes.NewReader(body)))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var o Owner
	_ = json.NewDecoder(rec.Body).Decode(&o)
	if o.Service != "payments" {
		t.Errorf("expected service=payments, got %q", o.Service)
	}
	if o.Team != "sre" {
		t.Errorf("expected team=sre, got %q", o.Team)
	}
}

func TestHandler_Get_KnownService(t *testing.T) {
	s := New()
	_ = s.Set(Owner{Service: "auth", Team: "identity", Email: "id@example.com"})
	h := NewHandler(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ownership/auth", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var o Owner
	_ = json.NewDecoder(rec.Body).Decode(&o)
	if o.Team != "identity" {
		t.Errorf("expected team=identity, got %q", o.Team)
	}
}

func TestHandler_Get_UnknownService(t *testing.T) {
	h := newHandler()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/ownership/ghost", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_Delete_RemovesOwner(t *testing.T) {
	s := New()
	_ = s.Set(Owner{Service: "svc", Team: "t", Email: "t@e.com"})
	h := NewHandler(s)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/ownership/svc", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h := newHandler()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/ownership/svc", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
