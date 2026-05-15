package healthscore

import (
	"testing"
)

func TestRecord_And_Get(t *testing.T) {
	s := New()
	if err := s.Record("api", 99.0, 50.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	sc, ok := s.Get("api")
	if !ok {
		t.Fatal("expected score to exist")
	}
	if sc.Service != "api" {
		t.Errorf("service = %q, want %q", sc.Service, "api")
	}
	if sc.Value <= 0 {
		t.Errorf("expected positive score, got %v", sc.Value)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected no score for unknown service")
	}
}

func TestRecord_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Record("", 100, 0); err == nil {
		t.Fatal("expected error for empty service")
	}
}

func TestRecord_InvalidUptimeReturnsError(t *testing.T) {
	s := New()
	if err := s.Record("api", 110, 0); err == nil {
		t.Fatal("expected error for uptime > 100")
	}
}

func TestGrade_A(t *testing.T) {
	s := New()
	_ = s.Record("api", 99.0, 10.0)
	sc, _ := s.Get("api")
	if sc.Grade != "A" {
		t.Errorf("grade = %q, want A", sc.Grade)
	}
}

func TestGrade_F_HighLatency(t *testing.T) {
	s := New()
	_ = s.Record("slow", 60.0, 3000.0)
	sc, _ := s.Get("slow")
	if sc.Grade != "F" {
		t.Errorf("grade = %q, want F", sc.Grade)
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	s := New()
	_ = s.Record("a", 99, 50)
	_ = s.Record("b", 80, 200)
	all := s.All()
	if len(all) != 2 {
		t.Errorf("len = %d, want 2", len(all))
	}
}

func TestRecord_Overwrites(t *testing.T) {
	s := New()
	_ = s.Record("api", 99, 50)
	_ = s.Record("api", 40, 100)
	sc, _ := s.Get("api")
	if sc.Grade == "A" {
		t.Error("expected grade to reflect updated (lower) score")
	}
}
