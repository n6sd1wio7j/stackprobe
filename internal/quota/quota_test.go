package quota

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	s := New(time.Minute)
	if err := s.Set("svc", 10); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, err := s.Get("svc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Limit != 10 {
		t.Errorf("expected limit 10, got %d", e.Limit)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(time.Minute)
	_, err := s.Get("missing")
	if err != ErrUnknownService {
		t.Errorf("expected ErrUnknownService, got %v", err)
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New(time.Minute)
	if err := s.Set("", 5); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestSet_ZeroLimitReturnsError(t *testing.T) {
	s := New(time.Minute)
	if err := s.Set("svc", 0); err == nil {
		t.Error("expected error for zero limit")
	}
}

func TestAllow_PermitsUpToLimit(t *testing.T) {
	s := New(time.Minute)
	_ = s.Set("svc", 3)
	for i := 0; i < 3; i++ {
		if err := s.Allow("svc"); err != nil {
			t.Fatalf("request %d should be allowed: %v", i+1, err)
		}
	}
	if err := s.Allow("svc"); err != ErrQuotaExceeded {
		t.Errorf("expected ErrQuotaExceeded, got %v", err)
	}
}

func TestAllow_UnknownService(t *testing.T) {
	s := New(time.Minute)
	if err := s.Allow("ghost"); err != ErrUnknownService {
		t.Errorf("expected ErrUnknownService, got %v", err)
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	s := New(50 * time.Millisecond)
	_ = s.Set("svc", 1)
	_ = s.Allow("svc")
	if err := s.Allow("svc"); err != ErrQuotaExceeded {
		t.Fatalf("expected quota exceeded before reset")
	}
	time.Sleep(60 * time.Millisecond)
	if err := s.Allow("svc"); err != nil {
		t.Errorf("expected allow after window reset, got %v", err)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New(time.Minute)
	_ = s.Set("a", 5)
	_ = s.Set("b", 10)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
