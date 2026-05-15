package lockout

import (
	"testing"
	"time"
)

func TestLock_And_IsLocked(t *testing.T) {
	s := New()
	if err := s.Lock("svc", "too many failures", time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !s.IsLocked("svc") {
		t.Fatal("expected service to be locked")
	}
}

func TestLock_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Lock("", "reason", time.Minute); err != ErrEmptyService {
		t.Fatalf("expected ErrEmptyService, got %v", err)
	}
}

func TestLock_AlreadyLockedReturnsError(t *testing.T) {
	s := New()
	_ = s.Lock("svc", "first", time.Minute)
	if err := s.Lock("svc", "second", time.Minute); err != ErrAlreadyLocked {
		t.Fatalf("expected ErrAlreadyLocked, got %v", err)
	}
}

func TestLock_ExpiredAllowsReLock(t *testing.T) {
	s := New()
	_ = s.Lock("svc", "old", -time.Second) // already expired
	if err := s.Lock("svc", "new", time.Minute); err != nil {
		t.Fatalf("unexpected error after expiry: %v", err)
	}
}

func TestUnlock_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Lock("svc", "reason", time.Minute)
	if err := s.Unlock("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.IsLocked("svc") {
		t.Fatal("expected service to be unlocked")
	}
}

func TestUnlock_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Unlock("ghost"); err != ErrNotLocked {
		t.Fatalf("expected ErrNotLocked, got %v", err)
	}
}

func TestGet_KnownService(t *testing.T) {
	s := New()
	_ = s.Lock("svc", "test reason", time.Minute)
	e, ok := s.Get("svc")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Reason != "test reason" {
		t.Fatalf("expected reason %q, got %q", "test reason", e.Reason)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected no entry for unknown service")
	}
}

func TestAll_ExcludesExpired(t *testing.T) {
	s := New()
	_ = s.Lock("active", "r", time.Minute)
	_ = s.Lock("expired", "r", -time.Second)
	all := s.All()
	if len(all) != 1 || all[0].Service != "active" {
		t.Fatalf("expected only active entry, got %v", all)
	}
}

func TestIsExpired_Fresh(t *testing.T) {
	e := Entry{ExpiresAt: time.Now().Add(time.Minute)}
	if e.IsExpired() {
		t.Fatal("expected entry to be fresh")
	}
}

func TestIsExpired_Stale(t *testing.T) {
	e := Entry{ExpiresAt: time.Now().Add(-time.Second)}
	if !e.IsExpired() {
		t.Fatal("expected entry to be expired")
	}
}
