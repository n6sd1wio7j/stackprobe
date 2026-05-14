package deadletter

import (
	"testing"
)

func TestPush_And_Get(t *testing.T) {
	s := New(10)
	if err := s.Push("svc-a", "timeout"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ents, err := s.Get("svc-a")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if len(ents) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(ents))
	}
	if ents[0].Reason != "timeout" {
		t.Errorf("expected reason 'timeout', got %q", ents[0].Reason)
	}
	if ents[0].OccurredAt.IsZero() {
		t.Error("expected OccurredAt to be set")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(10)
	_, err := s.Get("ghost")
	if err != ErrUnknownService {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestPush_EnforcesCapacity(t *testing.T) {
	s := New(2)
	_ = s.Push("svc", "err1")
	_ = s.Push("svc", "err2")
	err := s.Push("svc", "err3")
	if err != ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
}

func TestPush_EmptyServiceReturnsError(t *testing.T) {
	s := New(10)
	if err := s.Push("", "reason"); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestFlush_ClearsEntries(t *testing.T) {
	s := New(10)
	_ = s.Push("svc", "r1")
	_ = s.Push("svc", "r2")
	ents, err := s.Flush("svc")
	if err != nil {
		t.Fatalf("Flush error: %v", err)
	}
	if len(ents) != 2 {
		t.Fatalf("expected 2 flushed entries, got %d", len(ents))
	}
	_, err = s.Get("svc")
	if err != ErrUnknownService {
		t.Error("expected service to be removed after flush")
	}
}

func TestFlush_Unknown(t *testing.T) {
	s := New(10)
	_, err := s.Flush("ghost")
	if err != ErrUnknownService {
		t.Fatalf("expected ErrUnknownService, got %v", err)
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	s := New(10)
	_ = s.Push("svc-a", "err")
	_ = s.Push("svc-b", "err")
	_ = s.Push("svc-b", "err2")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 services, got %d", len(all))
	}
	if len(all["svc-b"]) != 2 {
		t.Errorf("expected 2 entries for svc-b, got %d", len(all["svc-b"]))
	}
}
