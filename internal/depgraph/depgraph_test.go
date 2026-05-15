package depgraph

import (
	"testing"
)

func TestRegister_And_Get(t *testing.T) {
	s := New()
	if err := s.Register("api", []string{"db", "cache"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, ok := s.Get("api")
	if !ok {
		t.Fatal("expected node to exist")
	}
	if n.Service != "api" {
		t.Errorf("expected api, got %s", n.Service)
	}
	if len(n.Children) != 2 {
		t.Errorf("expected 2 children, got %d", len(n.Children))
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestRegister_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Register("", nil); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSetHealth_KnownService(t *testing.T) {
	s := New()
	_ = s.Register("svc", nil)
	if err := s.SetHealth("svc", false); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	n, _ := s.Get("svc")
	if n.Healthy {
		t.Error("expected healthy=false")
	}
}

func TestSetHealth_UnknownService(t *testing.T) {
	s := New()
	if err := s.SetHealth("ghost", false); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestUnhealthy_ReturnsOnlyUnhealthy(t *testing.T) {
	s := New()
	_ = s.Register("a", nil)
	_ = s.Register("b", nil)
	_ = s.Register("c", nil)
	_ = s.SetHealth("b", false)
	list := s.Unhealthy()
	if len(list) != 1 || list[0] != "b" {
		t.Errorf("expected [b], got %v", list)
	}
}

func TestAll_ReturnsAllNodes(t *testing.T) {
	s := New()
	_ = s.Register("x", nil)
	_ = s.Register("y", []string{"x"})
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 nodes, got %d", len(all))
	}
}
