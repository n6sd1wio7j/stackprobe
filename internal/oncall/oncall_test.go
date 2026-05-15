package oncall

import (
	"testing"
	"time"
)

func futureTime(d time.Duration) time.Time {
	return time.Now().Add(d)
}

func TestSet_And_Get(t *testing.T) {
	s := New()
	until := futureTime(time.Hour)
	if err := s.Set("api", "alice", "+1000", "alice@example.com", until); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sc, ok := s.Get("api")
	if !ok {
		t.Fatal("expected schedule to exist")
	}
	if sc.Owner != "alice" {
		t.Errorf("expected alice, got %s", sc.Owner)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", "alice", "", "", futureTime(time.Hour)); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_EmptyOwnerReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("api", "", "", "", futureTime(time.Hour)); err == nil {
		t.Fatal("expected error for empty owner")
	}
}

func TestSet_PastUntilReturnsError(t *testing.T) {
	s := New()
	past := time.Now().Add(-time.Minute)
	if err := s.Set("api", "alice", "", "", past); err == nil {
		t.Fatal("expected error for past until")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("api", "alice", "", "", futureTime(time.Hour))
	if err := s.Delete("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("api")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestActive_ReturnsOnlyFuture(t *testing.T) {
	s := New()
	_ = s.Set("api", "alice", "", "", futureTime(time.Hour))
	// Manually insert an expired entry to simulate expiry.
	s.mu.Lock()
	s.records["old"] = Schedule{Service: "old", Owner: "bob", Until: time.Now().Add(-time.Minute)}
	s.mu.Unlock()

	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active schedule, got %d", len(active))
	}
	if active[0].Service != "api" {
		t.Errorf("expected api, got %s", active[0].Service)
	}
}

func TestAll_ReturnsAll(t *testing.T) {
	s := New()
	_ = s.Set("api", "alice", "", "", futureTime(time.Hour))
	_ = s.Set("db", "bob", "", "", futureTime(2*time.Hour))
	if len(s.All()) != 2 {
		t.Fatalf("expected 2 schedules")
	}
}

func TestSet_OverwritesExistingEntry(t *testing.T) {
	s := New()
	_ = s.Set("api", "alice", "", "", futureTime(time.Hour))
	if err := s.Set("api", "bob", "", "", futureTime(2*time.Hour)); err != nil {
		t.Fatalf("unexpected error on overwrite: %v", err)
	}
	sc, ok := s.Get("api")
	if !ok {
		t.Fatal("expected schedule to exist after overwrite")
	}
	if sc.Owner != "bob" {
		t.Errorf("expected bob after overwrite, got %s", sc.Owner)
	}
}
