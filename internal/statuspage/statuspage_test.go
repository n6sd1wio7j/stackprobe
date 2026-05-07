package statuspage

import (
	"testing"
)

func TestAdd_CreatesIncident(t *testing.T) {
	s := New()
	id := s.Add("api", "API Outage", "API is unreachable", SeverityMajor)
	if id == "" {
		t.Fatal("expected non-empty ID")
	}
	all := s.All()
	if len(all) != 1 {
		t.Fatalf("expected 1 incident, got %d", len(all))
	}
	if all[0].Service != "api" {
		t.Errorf("expected service 'api', got %s", all[0].Service)
	}
	if all[0].Severity != SeverityMajor {
		t.Errorf("expected severity major, got %s", all[0].Severity)
	}
}

func TestResolve_MarksResolved(t *testing.T) {
	s := New()
	id := s.Add("db", "DB slow", "High latency", SeverityMinor)
	ok := s.Resolve(id)
	if !ok {
		t.Fatal("expected Resolve to return true")
	}
	all := s.All()
	if all[0].ResolvedAt == nil {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestResolve_UnknownID(t *testing.T) {
	s := New()
	if s.Resolve("INC-9999") {
		t.Error("expected false for unknown ID")
	}
}

func TestResolve_AlreadyResolved(t *testing.T) {
	s := New()
	id := s.Add("svc", "t", "m", SeverityNone)
	s.Resolve(id)
	if s.Resolve(id) {
		t.Error("expected false when resolving an already-resolved incident")
	}
}

func TestActive_ReturnsOnlyUnresolved(t *testing.T) {
	s := New()
	id1 := s.Add("a", "t1", "m1", SeverityMinor)
	s.Add("b", "t2", "m2", SeverityCritical)
	s.Resolve(id1)
	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active incident, got %d", len(active))
	}
	if active[0].Service != "b" {
		t.Errorf("expected service 'b', got %s", active[0].Service)
	}
}

func TestAdd_SequentialIDs(t *testing.T) {
	s := New()
	id1 := s.Add("x", "t", "m", SeverityNone)
	id2 := s.Add("x", "t", "m", SeverityNone)
	if id1 == id2 {
		t.Error("expected unique IDs for successive incidents")
	}
}
