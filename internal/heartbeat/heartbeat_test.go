package heartbeat

import (
	"testing"
	"time"
)

func TestPing_SetsLastSeen(t *testing.T) {
	s := New(time.Minute)
	before := time.Now()
	s.Ping("svc-a")
	b, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected beat to exist")
	}
	if b.LastSeen.Before(before) {
		t.Error("LastSeen should be >= before ping")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(time.Minute)
	_, ok := s.Get("ghost")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestIsStale_Fresh(t *testing.T) {
	s := New(time.Minute)
	s.Ping("svc-b")
	b, _ := s.Get("svc-b")
	if b.IsStale(time.Now()) {
		t.Error("freshly pinged beat should not be stale")
	}
}

func TestIsStale_Old(t *testing.T) {
	s := New(time.Second)
	s.Ping("svc-c")
	b, _ := s.Get("svc-c")
	// Simulate old LastSeen
	b.LastSeen = time.Now().Add(-2 * time.Second)
	if !b.IsStale(time.Now()) {
		t.Error("old beat should be stale")
	}
}

func TestRegister_CustomInterval(t *testing.T) {
	s := New(time.Minute)
	s.Register("svc-d", 5*time.Second)
	b, ok := s.Get("svc-d")
	if !ok {
		t.Fatal("expected registered service to exist")
	}
	if b.Interval != 5*time.Second {
		t.Errorf("expected 5s interval, got %v", b.Interval)
	}
}

func TestStale_ReturnsOnlyStale(t *testing.T) {
	s := New(time.Minute)
	s.Ping("fresh")
	s.Ping("stale")

	// Make "stale" appear old by direct manipulation via Register then Ping trick
	// We'll register with tiny interval so it goes stale immediately
	s.Register("stale", time.Nanosecond)
	s.Ping("stale")
	time.Sleep(time.Millisecond)

	stales := s.Stale()
	found := false
	for _, b := range stales {
		if b.Service == "stale" {
			found = true
		}
		if b.Service == "fresh" {
			t.Error("fresh service should not appear in stale list")
		}
	}
	if !found {
		t.Error("expected stale service in stale list")
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	s := New(time.Minute)
	s.Ping("x")
	s.Ping("y")
	s.Ping("z")
	all := s.All()
	if len(all) != 3 {
		t.Errorf("expected 3 beats, got %d", len(all))
	}
}
