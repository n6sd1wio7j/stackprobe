package ownership

import (
	"testing"
)

func validOwner(service string) Owner {
	return Owner{Service: service, Team: "platform", Email: "platform@example.com", Slack: "#platform"}
}

func TestSet_And_Get(t *testing.T) {
	s := New()
	o := validOwner("api")
	if err := s.Set(o); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := s.Get("api")
	if !ok {
		t.Fatal("expected owner to be found")
	}
	if got.Team != "platform" {
		t.Errorf("got team %q, want %q", got.Team, "platform")
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected not found")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	err := s.Set(Owner{Team: "a", Email: "a@b.com"})
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_EmptyTeamReturnsError(t *testing.T) {
	s := New()
	err := s.Set(Owner{Service: "svc", Email: "a@b.com"})
	if err == nil {
		t.Fatal("expected error for empty team")
	}
}

func TestSet_EmptyEmailReturnsError(t *testing.T) {
	s := New()
	err := s.Set(Owner{Service: "svc", Team: "team"})
	if err == nil {
		t.Fatal("expected error for empty email")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New()
	_ = s.Set(validOwner("svc"))
	updated := Owner{Service: "svc", Team: "infra", Email: "infra@example.com"}
	_ = s.Set(updated)
	got, _ := s.Get("svc")
	if got.Team != "infra" {
		t.Errorf("expected updated team, got %q", got.Team)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set(validOwner("svc"))
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("svc")
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

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Set(validOwner("a"))
	_ = s.Set(validOwner("b"))
	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
