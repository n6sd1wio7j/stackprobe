package sla

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	target := Target{Service: "api", UptimePercent: 99.9, MaxLatencyMs: 200}
	s.Set(target)

	got, ok := s.Get("api")
	if !ok {
		t.Fatal("expected target to exist")
	}
	if got.UptimePercent != 99.9 {
		t.Errorf("expected 99.9, got %v", got.UptimePercent)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected no target for unknown service")
	}
}

func TestEvaluate_NoViolation(t *testing.T) {
	s := New()
	s.Set(Target{Service: "api", UptimePercent: 99.0, MaxLatencyMs: 300})
	s.Evaluate("api", 99.5, 150)

	if len(s.Violations()) != 0 {
		t.Fatalf("expected no violations, got %d", len(s.Violations()))
	}
}

func TestEvaluate_UptimeViolation(t *testing.T) {
	s := New()
	s.Set(Target{Service: "api", UptimePercent: 99.9, MaxLatencyMs: 0})
	s.Evaluate("api", 98.0, 0)

	vs := s.Violations()
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Kind != "uptime" {
		t.Errorf("expected uptime violation, got %s", vs[0].Kind)
	}
}

func TestEvaluate_LatencyViolation(t *testing.T) {
	s := New()
	s.Set(Target{Service: "api", UptimePercent: 99.0, MaxLatencyMs: 100})
	s.Evaluate("api", 99.5, 250)

	vs := s.Violations()
	if len(vs) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(vs))
	}
	if vs[0].Kind != "latency" {
		t.Errorf("expected latency violation, got %s", vs[0].Kind)
	}
}

func TestEvaluate_BothViolations(t *testing.T) {
	s := New()
	s.Set(Target{Service: "db", UptimePercent: 99.9, MaxLatencyMs: 50})
	s.Evaluate("db", 95.0, 200)

	if len(s.Violations()) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(s.Violations()))
	}
}

func TestEvaluate_UnknownServiceIsNoop(t *testing.T) {
	s := New()
	s.Evaluate("ghost", 0, 9999)
	if len(s.Violations()) != 0 {
		t.Fatal("expected no violations for unregistered service")
	}
}

func TestTargets_ReturnsAll(t *testing.T) {
	s := New()
	s.Set(Target{Service: "a", UptimePercent: 99.0})
	s.Set(Target{Service: "b", UptimePercent: 95.0})

	if len(s.Targets()) != 2 {
		t.Fatalf("expected 2 targets, got %d", len(s.Targets()))
	}
}
