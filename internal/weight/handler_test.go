package weight

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newHandler(t *testing.T) (http.Handler, *Store) {
	t.Helper()
	s := New(1)
	return NewHandler(s), s
}

func TestHandler_GetAll_Empty(t *testing.T) {
	h, _ := newHandler(t)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/weights", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestHandler_Put_SetsWeight(t *testing.T) {
	h, s := newHandler(t)
	body, _ := json.Marshal(map[string]int{"weight": 8})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/weights/api", bytes.NewReader(body)))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if s.Get("api") != 8 {
		t.Fatalf("expected weight 8")
	}
}

func TestHandler_Put_InvalidWeight(t *testing.T) {
	h, _ := newHandler(t)
	body, _ := json.Marshal(map[string]int{"weight": 0})
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPut, "/weights/api", bytes.NewReader(body)))
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestHandler_Delete_RemovesWeight(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("svc", 5)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodDelete, "/weights/svc", nil))
	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}
	if s.Get("svc") != 1 {
		t.Fatalf("expected revert to default 1")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	h, _ := newHandler(t)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/weights", nil))
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestHandler_Get_SingleService(t *testing.T) {
	h, s := newHandler(t)
	_ = s.Set("db", 3)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/weights/db", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var resp map[string]int
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp["weight"] != 3 {
		t.Fatalf("expected 3, got %d", resp["weight"])
	}
}
