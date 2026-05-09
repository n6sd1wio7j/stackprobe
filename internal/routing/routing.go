package routing

import (
	"errors"
	"net/http"
	"strings"
	"sync"
)

// Rule maps a path prefix to a named upstream target.
type Rule struct {
	Prefix   string `json:"prefix"`
	Target   string `json:"target"`
	StripPrefix bool `json:"strip_prefix"`
}

// Router matches incoming requests to upstream targets based on prefix rules.
type Router struct {
	mu    sync.RWMutex
	rules []Rule
}

// New returns an empty Router.
func New() *Router {
	return &Router{}
}

// Add registers a new routing rule. Returns an error if prefix is empty or
// target is empty.
func (r *Router) Add(rule Rule) error {
	if strings.TrimSpace(rule.Prefix) == "" {
		return errors.New("routing: prefix must not be empty")
	}
	if strings.TrimSpace(rule.Target) == "" {
		return errors.New("routing: target must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rules = append(r.rules, rule)
	return nil
}

// Remove deletes all rules whose prefix matches the given value.
func (r *Router) Remove(prefix string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	filtered := r.rules[:0]
	for _, rule := range r.rules {
		if rule.Prefix != prefix {
			filtered = append(filtered, rule)
		}
	}
	r.rules = filtered
}

// Match returns the first Rule whose prefix matches the request path.
// The second return value indicates whether a match was found.
func (r *Router) Match(req *http.Request) (Rule, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	path := req.URL.Path
	for _, rule := range r.rules {
		if strings.HasPrefix(path, rule.Prefix) {
			return rule, true
		}
	}
	return Rule{}, false
}

// All returns a snapshot of all registered rules.
func (r *Router) All() []Rule {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Rule, len(r.rules))
	copy(out, r.rules)
	return out
}
