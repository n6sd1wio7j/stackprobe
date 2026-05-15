package incident

import (
	"testing"
)

func TestOpen_And_Get(t *testing.T) {
	s := New()
	id, err := s.Open("api", "High error rate", SeverityHigh)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inc, ok := s.Get(id)
	if !ok {
		t.Fatal("expected incident to exist")
	}
	if inc.Service != "api" || inc.Title != "High error rate" {
		t.Errorf("unexpected incident fields: %+v", inc)
	}
	if inc.Resolved {
		t.Error("new incident should not be resolved")
	}
}

func TestOpen_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	_, err := s.Open("", "title", SeverityLow)
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestOpen_EmptyTitleReturnsError(t *testing.T) {
	s := New()
	_, err := s.Open("svc", "", SeverityLow)
	if err == nil {
		t.Fatal("expected error for empty title")
	}
}

func TestResolve_MarksResolved(t *testing.T) {
	s := New()
	id, _ := s.Open("db", "Connection pool exhausted", SeverityCritical)
	if err := s.Resolve(id); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	inc, _ := s.Get(id)
	if !inc.Resolved {
		t.Error("expected incident to be resolved")
	}
	if inc.ResolvedAt == nil {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestResolve_UnknownID(t *testing.T) {
	s := New()
	if err := s.Resolve("INC-9999"); err == nil {
		t.Fatal("expected error for unknown ID")
	}
}

func TestResolve_AlreadyResolved(t *testing.T) {
	s := New()
	id, _ := s.Open("svc", "title", SeverityLow)
	s.Resolve(id)
	if err := s.Resolve(id); err == nil {
		t.Fatal("expected error when resolving already-resolved incident")
	}
}

func TestActive_ReturnsOnlyUnresolved(t *testing.T) {
	s := New()
	id1, _ := s.Open("svc-a", "incident one", SeverityHigh)
	s.Open("svc-b", "incident two", SeverityLow)
	s.Resolve(id1)

	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active incident, got %d", len(active))
	}
	if active[0].Service != "svc-b" {
		t.Errorf("unexpected active incident: %+v", active[0])
	}
}

func TestAll_ReturnsBothStates(t *testing.T) {
	s := New()
	id1, _ := s.Open("svc", "first", SeverityCritical)
	s.Open("svc", "second", SeverityLow)
	s.Resolve(id1)

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 incidents, got %d", len(all))
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("INC-0000")
	if ok {
		t.Fatal("expected not found for unknown ID")
	}
}
