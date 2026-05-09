package audit

import (
	"testing"
)

func TestRecord_StoresEvent(t *testing.T) {
	s := New(10)
	e := s.Record(KindCheck, "api", "scheduler", "check completed")
	if e.ID == "" {
		t.Fatal("expected non-empty ID")
	}
	if e.Kind != KindCheck {
		t.Fatalf("expected kind %q, got %q", KindCheck, e.Kind)
	}
	if e.Service != "api" {
		t.Fatalf("expected service api, got %q", e.Service)
	}
	if e.Timestamp.IsZero() {
		t.Fatal("expected non-zero timestamp")
	}
}

func TestAll_NewestFirst(t *testing.T) {
	s := New(10)
	s.Record(KindCheck, "svc-a", "sys", "first")
	s.Record(KindAlert, "svc-b", "sys", "second")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 events, got %d", len(all))
	}
	if all[0].Message != "second" {
		t.Fatalf("expected newest first, got %q", all[0].Message)
	}
}

func TestRecord_EnforcesLimit(t *testing.T) {
	s := New(3)
	for i := 0; i < 5; i++ {
		s.Record(KindConfig, "svc", "user", "msg")
	}
	if len(s.All()) != 3 {
		t.Fatalf("expected 3 events after overflow, got %d", len(s.All()))
	}
}

func TestFilter_ByKind(t *testing.T) {
	s := New(20)
	s.Record(KindCheck, "svc", "sys", "check")
	s.Record(KindAlert, "svc", "sys", "alert")
	s.Record(KindCheck, "svc", "sys", "check2")

	checks := s.Filter(KindCheck)
	if len(checks) != 2 {
		t.Fatalf("expected 2 check events, got %d", len(checks))
	}
	for _, e := range checks {
		if e.Kind != KindCheck {
			t.Fatalf("unexpected kind %q in filter result", e.Kind)
		}
	}
}

func TestFilter_EmptyKindReturnsAll(t *testing.T) {
	s := New(20)
	s.Record(KindCheck, "svc", "sys", "a")
	s.Record(KindAlert, "svc", "sys", "b")
	if len(s.Filter("")) != 2 {
		t.Fatal("expected all events when kind is empty")
	}
}

func TestRecord_SequentialIDs(t *testing.T) {
	s := New(10)
	e1 := s.Record(KindCheck, "svc", "sys", "first")
	e2 := s.Record(KindCheck, "svc", "sys", "second")
	if e1.ID == e2.ID {
		t.Fatal("expected unique IDs for sequential events")
	}
}
