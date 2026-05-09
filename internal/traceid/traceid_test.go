package traceid

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate_UniqueIDs(t *testing.T) {
	a, err := Generate()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, _ := Generate()
	if a == b {
		t.Error("expected unique IDs, got duplicates")
	}
	if len(a) != 32 {
		t.Errorf("expected 32 hex chars, got %d", len(a))
	}
}

func TestStore_SetAndGet(t *testing.T) {
	s := New()
	s.Set("svc-a", "abc123")
	v, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected key to be found")
	}
	if v != "abc123" {
		t.Errorf("expected abc123, got %s", v)
	}
}

func TestStore_Get_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Error("expected not found for unknown key")
	}
}

func TestStore_Delete(t *testing.T) {
	s := New()
	s.Set("svc", "id1")
	s.Delete("svc")
	_, ok := s.Get("svc")
	if ok {
		t.Error("expected key to be deleted")
	}
}

func TestMiddleware_GeneratesID(t *testing.T) {
	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := FromContext(r.Context())
		if id == "" {
			t.Error("expected trace ID in context")
		}
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	handler.ServeHTTP(rec, req)

	if rec.Header().Get(Header) == "" {
		t.Error("expected X-Trace-Id response header")
	}
}

func TestMiddleware_PropagatesExistingID(t *testing.T) {
	const existing = "deadbeefdeadbeefdeadbeefdeadbeef"
	var captured string

	handler := Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		captured = FromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	}))

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(Header, existing)
	handler.ServeHTTP(rec, req)

	if captured != existing {
		t.Errorf("expected %s, got %s", existing, captured)
	}
	if rec.Header().Get(Header) != existing {
		t.Error("expected propagated ID in response header")
	}
}

func TestFromContext_Empty(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	id := FromContext(req.Context())
	if id != "" {
		t.Errorf("expected empty string, got %s", id)
	}
}
