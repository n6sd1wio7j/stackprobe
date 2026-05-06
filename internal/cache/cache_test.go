package cache

import (
	"testing"
	"time"
)

func TestSet_And_Get(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("svc-a", "up", 42)

	e, ok := c.Get("svc-a")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Status != "up" {
		t.Errorf("status: got %q, want %q", e.Status, "up")
	}
	if e.LatencyMs != 42 {
		t.Errorf("latency: got %d, want 42", e.LatencyMs)
	}
}

func TestGet_Unknown(t *testing.T) {
	c := New(5 * time.Second)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for unknown service")
	}
}

func TestIsExpired_Fresh(t *testing.T) {
	c := New(10 * time.Second)
	c.Set("svc", "up", 0)
	e, _ := c.Get("svc")
	if e.IsExpired() {
		t.Error("fresh entry should not be expired")
	}
}

func TestIsExpired_Stale(t *testing.T) {
	c := New(-1 * time.Second) // negative TTL → already expired
	c.Set("svc", "up", 0)
	e, _ := c.Get("svc")
	if !e.IsExpired() {
		t.Error("entry with negative TTL should be expired")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("svc", "up", 10)
	c.Delete("svc")
	_, ok := c.Get("svc")
	if ok {
		t.Fatal("expected entry to be removed after Delete")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("alpha", "up", 1)
	c.Set("beta", "down", 2)

	all := c.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["alpha"].Status != "up" {
		t.Errorf("alpha status: got %q", all["alpha"].Status)
	}
	if all["beta"].Status != "down" {
		t.Errorf("beta status: got %q", all["beta"].Status)
	}
}

func TestSet_Overwrites(t *testing.T) {
	c := New(5 * time.Second)
	c.Set("svc", "up", 10)
	c.Set("svc", "down", 99)
	e, _ := c.Get("svc")
	if e.Status != "down" {
		t.Errorf("expected overwritten status %q, got %q", "down", e.Status)
	}
	if e.LatencyMs != 99 {
		t.Errorf("expected overwritten latency 99, got %d", e.LatencyMs)
	}
}
