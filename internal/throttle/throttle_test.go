package throttle

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAllow_PermitsUpToLimit(t *testing.T) {
	s := New(3, time.Minute)
	for i := 0; i < 3; i++ {
		if !s.Allow("svc") {
			t.Fatalf("expected Allow to return true on call %d", i+1)
		}
	}
	if s.Allow("svc") {
		t.Fatal("expected Allow to return false after limit reached")
	}
}

func TestAllow_IndependentServices(t *testing.T) {
	s := New(1, time.Minute)
	if !s.Allow("a") {
		t.Fatal("expected first call for 'a' to be allowed")
	}
	if !s.Allow("b") {
		t.Fatal("expected first call for 'b' to be allowed")
	}
	if s.Allow("a") {
		t.Fatal("expected second call for 'a' to be blocked")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	s := New(1, 50*time.Millisecond)
	if !s.Allow("svc") {
		t.Fatal("expected first call to be allowed")
	}
	if s.Allow("svc") {
		t.Fatal("expected second call to be blocked")
	}
	time.Sleep(60 * time.Millisecond)
	if !s.Allow("svc") {
		t.Fatal("expected call after window expiry to be allowed")
	}
}

func TestRemaining_DecreasesWithUse(t *testing.T) {
	s := New(5, time.Minute)
	if got := s.Remaining("svc"); got != 5 {
		t.Fatalf("expected 5 remaining, got %d", got)
	}
	s.Allow("svc")
	s.Allow("svc")
	if got := s.Remaining("svc"); got != 3 {
		t.Fatalf("expected 3 remaining, got %d", got)
	}
}

func TestReset_AllowsImmediateRequest(t *testing.T) {
	s := New(1, time.Minute)
	s.Allow("svc")
	if s.Allow("svc") {
		t.Fatal("expected blocked before reset")
	}
	s.Reset("svc")
	if !s.Allow("svc") {
		t.Fatal("expected allowed after reset")
	}
}

func TestMiddleware_ThrottlesAt429(t *testing.T) {
	s := New(2, time.Minute)
	handler := s.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	makeReq := func(service string) int {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		if service != "" {
			req.Header.Set("X-Service", service)
		}
		handler.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := makeReq("api"); code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}
	if code := makeReq("api"); code != http.StatusOK {
		t.Fatalf("expected 200, got %d", code)
	}
	if code := makeReq("api"); code != http.StatusTooManyRequests {
		t.Fatalf("expected 429, got %d", code)
	}
}
