package semver

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	v, err := Parse("1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("got %+v", v)
	}
}

func TestParse_VPrefix(t *testing.T) {
	v, err := Parse("v2.0.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 2 || v.Minor != 0 || v.Patch != 1 {
		t.Fatalf("got %+v", v)
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []string{"", "1.2", "a.b.c", "1.x.3"}
	for _, c := range cases {
		_, err := Parse(c)
		if err == nil {
			t.Errorf("expected error for %q", c)
		}
	}
}

func TestCompare(t *testing.T) {
	parse := func(s string) Version { v, _ := Parse(s); return v }
	tests := []struct {
		a, b string
		want int
	}{
		{"1.0.0", "1.0.0", 0},
		{"1.0.0", "2.0.0", -1},
		{"2.0.0", "1.0.0", 1},
		{"1.2.3", "1.2.4", -1},
		{"1.3.0", "1.2.9", 1},
	}
	for _, tt := range tests {
		got := Compare(parse(tt.a), parse(tt.b))
		if got != tt.want {
			t.Errorf("Compare(%s, %s) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestSet_And_Get(t *testing.T) {
	s := New()
	if err := s.Set("api", "1.2.3"); err != nil {
		t.Fatalf("Set: %v", err)
	}
	v, ok := s.Get("api")
	if !ok {
		t.Fatal("expected entry")
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Fatalf("got %+v", v)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := New()
	_, ok := s.Get("missing")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestSet_EmptyServiceReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("", "1.0.0"); err == nil {
		t.Fatal("expected error")
	}
}

func TestSet_InvalidVersionReturnsError(t *testing.T) {
	s := New()
	if err := s.Set("api", "not-a-version"); err == nil {
		t.Fatal("expected error")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := New()
	_ = s.Set("api", "1.0.0")
	s.Delete("api")
	_, ok := s.Get("api")
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	s := New()
	_ = s.Set("svc-a", "1.0.0")
	_ = s.Set("svc-b", "2.1.0")
	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(all))
	}
}
