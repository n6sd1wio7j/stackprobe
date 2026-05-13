package priority

import (
	"testing"
)

func TestSet_And_Get(t *testing.T) {
	s := New(Low)
	if err := s.Set("svc-a", High); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("svc-a"); got != High {
		t.Fatalf("expected High, got %v", got)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New(Medium)
	if got := s.Get("unknown"); got != Medium {
		t.Fatalf("expected default Medium, got %v", got)
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New(Low)
	if err := s.Set("", High); err == nil {
		t.Fatal("expected error for empty service name")
	}
}

func TestSet_InvalidLevelReturnsError(t *testing.T) {
	s := New(Low)
	if err := s.Set("svc", Level(99)); err == nil {
		t.Fatal("expected error for invalid level")
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := New(Low)
	_ = s.Set("svc", Low)
	_ = s.Set("svc", Critical)
	if got := s.Get("svc"); got != Critical {
		t.Fatalf("expected Critical after overwrite, got %v", got)
	}
}

func TestDelete_RevertsToDefault(t *testing.T) {
	s := New(Medium)
	_ = s.Set("svc", Critical)
	if err := s.Delete("svc"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.Get("svc"); got != Medium {
		t.Fatalf("expected default Medium after delete, got %v", got)
	}
}

func TestDelete_UnknownReturnsError(t *testing.T) {
	s := New(Low)
	if err := s.Delete("ghost"); err == nil {
		t.Fatal("expected error deleting unknown service")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New(Low)
	_ = s.Set("a", High)
	_ = s.Set("b", Critical)
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
	if all["a"] != High {
		t.Errorf("expected High for 'a'")
	}
}

func TestFilter_MinLevel(t *testing.T) {
	s := New(Low)
	_ = s.Set("low-svc", Low)
	_ = s.Set("med-svc", Medium)
	_ = s.Set("high-svc", High)
	_ = s.Set("crit-svc", Critical)

	result := s.Filter(High)
	if len(result) != 2 {
		t.Fatalf("expected 2 services at or above High, got %d", len(result))
	}
	found := map[string]bool{}
	for _, svc := range result {
		found[svc] = true
	}
	if !found["high-svc"] || !found["crit-svc"] {
		t.Errorf("expected high-svc and crit-svc in result")
	}
}

func TestLevelString(t *testing.T) {
	if Critical.String() != "critical" {
		t.Errorf("expected 'critical', got %q", Critical.String())
	}
	if Level(99).String() != "unknown" {
		t.Errorf("expected 'unknown' for invalid level")
	}
}
