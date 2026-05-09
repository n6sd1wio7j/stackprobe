package routing

import (
	"net/http"
	"testing"
)

func newRequest(path string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, "http://localhost"+path, nil)
	return req
}

func TestAdd_And_Match(t *testing.T) {
	r := New()
	if err := r.Add(Rule{Prefix: "/api", Target: "backend-api"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rule, ok := r.Match(newRequest("/api/users"))
	if !ok {
		t.Fatal("expected a match")
	}
	if rule.Target != "backend-api" {
		t.Errorf("expected backend-api, got %s", rule.Target)
	}
}

func TestMatch_NoMatch(t *testing.T) {
	r := New()
	_ = r.Add(Rule{Prefix: "/api", Target: "backend-api"})
	_, ok := r.Match(newRequest("/health"))
	if ok {
		t.Fatal("expected no match")
	}
}

func TestAdd_EmptyPrefix(t *testing.T) {
	r := New()
	if err := r.Add(Rule{Prefix: "", Target: "svc"}); err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestAdd_EmptyTarget(t *testing.T) {
	r := New()
	if err := r.Add(Rule{Prefix: "/x", Target: ""}); err == nil {
		t.Fatal("expected error for empty target")
	}
}

func TestRemove_DeletesRule(t *testing.T) {
	r := New()
	_ = r.Add(Rule{Prefix: "/api", Target: "backend-api"})
	_ = r.Add(Rule{Prefix: "/web", Target: "frontend"})
	r.Remove("/api")
	_, ok := r.Match(newRequest("/api/v1"))
	if ok {
		t.Fatal("expected no match after removal")
	}
	_, ok = r.Match(newRequest("/web/index"))
	if !ok {
		t.Fatal("expected /web rule to still exist")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	r := New()
	_ = r.Add(Rule{Prefix: "/a", Target: "svc-a"})
	_ = r.Add(Rule{Prefix: "/b", Target: "svc-b"})
	rules := r.All()
	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}
}

func TestMatch_StripPrefixFlag(t *testing.T) {
	r := New()
	_ = r.Add(Rule{Prefix: "/svc", Target: "backend", StripPrefix: true})
	rule, ok := r.Match(newRequest("/svc/ping"))
	if !ok {
		t.Fatal("expected match")
	}
	if !rule.StripPrefix {
		t.Error("expected StripPrefix to be true")
	}
}
