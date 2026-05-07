package timeout_test

import (
	"testing"
	"time"

	"github.com/your-org/stackprobe/internal/timeout"
)

func TestGet_ReturnsDefault(t *testing.T) {
	s := timeout.New(5 * time.Second)
	got := s.Get("svc-a")
	if got != 5*time.Second {
		t.Fatalf("expected 5s default, got %v", got)
	}
}

func TestSet_And_Get(t *testing.T) {
	s := timeout.New(5 * time.Second)
	s.Set("svc-a", 2*time.Second)
	got := s.Get("svc-a")
	if got != 2*time.Second {
		t.Fatalf("expected 2s override, got %v", got)
	}
}

func TestSet_DoesNotAffectOtherServices(t *testing.T) {
	s := timeout.New(10 * time.Second)
	s.Set("svc-a", 1*time.Second)
	got := s.Get("svc-b")
	if got != 10*time.Second {
		t.Fatalf("expected default for svc-b, got %v", got)
	}
}

func TestDelete_RevertsToDefault(t *testing.T) {
	s := timeout.New(7 * time.Second)
	s.Set("svc-a", 1*time.Second)
	s.Delete("svc-a")
	got := s.Get("svc-a")
	if got != 7*time.Second {
		t.Fatalf("expected default after delete, got %v", got)
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := timeout.New(5 * time.Second)
	s.Set("svc-a", 1*time.Second)
	s.Set("svc-b", 3*time.Second)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 overrides, got %d", len(all))
	}
	if all["svc-a"] != 1*time.Second {
		t.Errorf("svc-a: expected 1s, got %v", all["svc-a"])
	}
	if all["svc-b"] != 3*time.Second {
		t.Errorf("svc-b: expected 3s, got %v", all["svc-b"])
	}
}

func TestAll_IsolatedFromInternalState(t *testing.T) {
	s := timeout.New(5 * time.Second)
	s.Set("svc-a", 2*time.Second)
	snap := s.All()
	snap["svc-a"] = 99 * time.Second // mutate snapshot
	if s.Get("svc-a") != 2*time.Second {
		t.Fatal("snapshot mutation affected store")
	}
}
