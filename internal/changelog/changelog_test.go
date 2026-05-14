package changelog

import (
	"testing"
)

func TestAdd_And_Get(t *testing.T) {
	s := New(100)
	e, err := s.Add("api", KindAdded, "initial release", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Service != "api" || e.Kind != KindAdded {
		t.Fatalf("unexpected entry: %+v", e)
	}
	entries := s.Get("api")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(100)
	if got := s.Get("ghost"); len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

func TestAdd_EmptyServiceReturnsError(t *testing.T) {
	s := New(100)
	_, err := s.Add("", KindFixed, "fix", "bob")
	if err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestAdd_EmptyMessageReturnsError(t *testing.T) {
	s := New(100)
	_, err := s.Add("svc", KindFixed, "", "bob")
	if err == nil {
		t.Fatal("expected error for empty message")
	}
}

func TestAdd_EnforcesLimit(t *testing.T) {
	s := New(3)
	for i := 0; i < 5; i++ {
		s.Add("svc", KindChanged, "msg", "")
	}
	if len(s.All()) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(s.All()))
	}
}

func TestAll_NewestFirst(t *testing.T) {
	s := New(100)
	s.Add("svc", KindAdded, "first", "")
	s.Add("svc", KindChanged, "second", "")
	all := s.All()
	if all[0].Message != "second" {
		t.Fatalf("expected newest first, got %q", all[0].Message)
	}
}

func TestFilter_ByKind(t *testing.T) {
	s := New(100)
	s.Add("svc", KindAdded, "add", "")
	s.Add("svc", KindFixed, "fix", "")
	s.Add("svc", KindFixed, "fix2", "")
	got := s.Filter(KindFixed)
	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestFilter_EmptyKindReturnsAll(t *testing.T) {
	s := New(100)
	s.Add("svc", KindAdded, "a", "")
	s.Add("svc", KindRemoved, "b", "")
	if len(s.Filter("")) != 2 {
		t.Fatal("expected all entries")
	}
}
