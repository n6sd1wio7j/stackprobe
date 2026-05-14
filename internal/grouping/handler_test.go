package grouping

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) (http.Handler, *Store) {
	t.Helper()
	s := New()
	return NewHandler(s), s
}

func TestHandler_GetAll_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/groups", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestHandler_Put_CreatesGroup(t *testing.T) {
	h, s := newHandler(t)
	body, _ := json.Marshal(map[string]any{"services": []string{"svc-a", "svc-b"}})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/groups/mygroup", bytes.NewReader(body)))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	svcs, ok := s.Get("mygroup")
	if !ok || len(svcs) != 2 {
		t.Fatal("group not stored correctly")
	}
}

func TestHandler_Put_InvalidJSON(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/groups/g", bytes.NewBufferString("{")))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHandler_Get_KnownGroup(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("team", []string{"x", "y"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/groups/team", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var g Group
	if err := json.NewDecoder(rec.Body).Decode(&g); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if g.Name != "team" || len(g.Services) != 2 {
		t.Fatalf("unexpected group: %+v", g)
	}
}

func TestHandler_Get_UnknownGroup(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/groups/nope", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestHandler_Delete_RemovesGroup(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("del", []string{"a"})
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/groups/del", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	_, ok := s.Get("del")
	if ok {
		t.Fatal("group should have been deleted")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/groups", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}
