package alerting

import (
	"testing"
)

func TestEvaluate_FiresAlertOnDown(t *testing.T) {
	m := New()
	alert := m.Evaluate("svc-a", false)
	if alert == nil {
		t.Fatal("expected alert to be fired")
	}
	if alert.Service != "svc-a" {
		t.Errorf("expected service svc-a, got %s", alert.Service)
	}
	if alert.ResolvedAt != nil {
		t.Error("expected ResolvedAt to be nil for a new alert")
	}
}

func TestEvaluate_NoAlertWhenAlreadyDown(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", false) // first down
	alert := m.Evaluate("svc-a", false) // still down
	if alert != nil {
		t.Error("expected no duplicate alert")
	}
}

func TestEvaluate_ResolvesAlert(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", false)
	resolved := m.Evaluate("svc-a", true)
	if resolved == nil {
		t.Fatal("expected resolved alert")
	}
	if resolved.ResolvedAt == nil {
		t.Error("expected ResolvedAt to be set")
	}
}

func TestEvaluate_NoAlertWhenAlreadyUp(t *testing.T) {
	m := New()
	alert := m.Evaluate("svc-a", true)
	if alert != nil {
		t.Error("expected no alert when service is healthy")
	}
}

func TestActive_ReturnsOnlyFiring(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", false)
	m.Evaluate("svc-b", false)
	m.Evaluate("svc-a", true) // resolve svc-a

	active := m.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active alert, got %d", len(active))
	}
	if active[0].Service != "svc-b" {
		t.Errorf("expected svc-b, got %s", active[0].Service)
	}
}

func TestActive_EmptyWhenAllResolved(t *testing.T) {
	m := New()
	m.Evaluate("svc-a", false)
	m.Evaluate("svc-a", true)
	if len(m.Active()) != 0 {
		t.Error("expected no active alerts after resolution")
	}
}

func TestEvaluate_FiredAtIsSet(t *testing.T) {
	m := New()
	alert := m.Evaluate("svc-a", false)
	if alert == nil {
		t.Fatal("expected alert to be fired")
	}
	if alert.FiredAt.IsZero() {
		t.Error("expected FiredAt to be set on new alert")
	}
}

func TestEvaluate_ResolvedAtAfterFiredAt(t *testing.T) {
	m := New()
	fired := m.Evaluate("svc-a", false)
	if fired == nil {
		t.Fatal("expected alert to be fired")
	}
	resolved := m.Evaluate("svc-a", true)
	if resolved == nil {
		t.Fatal("expected resolved alert")
	}
	if resolved.ResolvedAt.Before(resolved.FiredAt) {
		t.Error("expected ResolvedAt to be after or equal to FiredAt")
	}
}
