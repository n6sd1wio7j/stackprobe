package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Route maps a path prefix to a backend URL.
type Route struct {
	Prefix  string
	Backend string
}

// Proxy holds a set of reverse-proxy routes.
type Proxy struct {
	mu      sync.RWMutex
	routes  []Route
	proxies map[string]*httputil.ReverseProxy
}

// New creates a Proxy pre-loaded with the given routes.
func New(routes []Route) (*Proxy, error) {
	p := &Proxy{
		proxies: make(map[string]*httputil.ReverseProxy),
	}
	for _, r := range routes {
		if err := p.Add(r); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Add registers a new route at runtime.
func (p *Proxy) Add(r Route) error {
	target, err := url.Parse(r.Backend)
	if err != nil {
		return err
	}
	rp := httputil.NewSingleHostReverseProxy(target)
	p.mu.Lock()
	defer p.mu.Unlock()
	p.routes = append(p.routes, r)
	p.proxies[r.Prefix] = rp
	return nil
}

// Routes returns a snapshot of current routes.
func (p *Proxy) Routes() []Route {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Route, len(p.routes))
	copy(out, p.routes)
	return out
}

// ServeHTTP matches the longest prefix and forwards the request.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	var best string
	for prefix := range p.proxies {
		if len(prefix) > len(best) && len(r.URL.Path) >= len(prefix) && r.URL.Path[:len(prefix)] == prefix {
			best = prefix
		}
	}
	if best == "" {
		http.Error(w, "no matching proxy route", http.StatusBadGateway)
		return
	}
	p.proxies[best].ServeHTTP(w, r)
}
