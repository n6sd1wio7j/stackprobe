package healthscore

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) (http.Handler, *Store) {
	t.Helper()
	store := New()
	return NewHandler(store), store
}

func TestHandler_GetAll_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthscore", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var out []Score
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out) != 0 {
		t.Errorf("expected empty list, got %d items", len(out))
	}
}

func TestHandler_GetAll_ReturnsScores(t *testing.T) {
	h, store := newHandler(t)
	_ = store.Record("api", 98, 100)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthscore", nil))
	var out []Score
	_ = json.NewDecoder(rec.Body).Decode(&out)
	if len(out) != 1 {
		t.Errorf("expected 1 score, got %d", len(out))
	}
}

func TestHandler_GetOne_KnownService(t *testing.T) {
	h, store := newHandler(t)
	_ = store.Record("db", 95, 200)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthscore/db", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var sc Score
	_ = json.NewDecoder(rec.Body).Decode(&sc)
	if sc.Service != "db" {
		t.Errorf("service = %q, want db", sc.Service)
	}
}

func TestHandler_GetOne_UnknownService(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthscore/ghost", nil))
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", rec.Code)
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthscore", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want 405", rec.Code)
	}
}
