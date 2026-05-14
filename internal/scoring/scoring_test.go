package scoring

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("api", 99.5, 45); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sc, ok := s.Get("api")
	if !ok {
		t.Fatal("expected score to exist")
	}
	if sc.Service != "api" {
		t.Errorf("got service %q, want %q", sc.Service, "api")
	}
	if sc.Grade != "A" {
		t.Errorf("got grade %q, want A", sc.Grade)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("ghost")
	if ok {
		t.Fatal("expected no score for unknown service")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", 100, 10); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestSet_InvalidUptimeReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("svc", 110, 10); err == nil {
		t.Fatal("expected error for uptime > 100")
	}
	if err := s.Set("svc", -1, 10); err == nil {
		t.Fatal("expected error for negative uptime")
	}
}

func TestGrade_LowUptime(t *testing.T) {
	s := New()
	_ = s.Set("bad", 50, 2000)
	sc, _ := s.Get("bad")
	if sc.Grade != "F" {
		t.Errorf("expected F, got %q", sc.Grade)
	}
}

func TestAll_ReturnsAllScores(t *testing.T) {
	s := New()
	_ = s.Set("a", 99, 30)
	_ = s.Set("b", 80, 200)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 scores, got %d", len(all))
	}
}

func TestDelete_RemovesScore(t *testing.T) {
	s := New()
	_ = s.Set("api", 99, 50)
	if err := s.Delete("api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := s.Get("api"); ok {
		t.Fatal("expected score to be deleted")
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New()
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error for unknown service")
	}
}

func TestReasoning_LowUptime(t *testing.T) {
	s := New()
	_ = s.Set("slow", 85, 100)
	sc, _ := s.Get("slow")
	if sc.Reasoning != "low uptime dragging score down" {
		t.Errorf("unexpected reasoning: %q", sc.Reasoning)
	}
}
