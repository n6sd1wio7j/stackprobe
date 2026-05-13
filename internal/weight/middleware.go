package weight

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
)

// WeightedBackend pairs a reverse proxy with its weight.
type WeightedBackend struct {
	Proxy  *httputil.ReverseProxy
	Weight int
}

// RoundRobinBalancer selects backends using weighted round-robin.
type RoundRobinBalancer struct {
	backends []WeightedBackend
	counter  atomic.Uint64
}

// NewBalancer builds a balancer from a map of target URL → service name,
// looking up each service weight from the store.
func NewBalancer(targets map[string]string, store *Store) (*RoundRobinBalancer, error) {
	var expanded []WeightedBackend
	for rawURL, svc := range targets {
		u, err := url.Parse(rawURL)
		if err != nil {
			return nil, err
		}
		w := store.Get(svc)
		proxy := httputil.NewSingleHostReverseProxy(u)
		for i := 0; i < w; i++ {
			expanded = append(expanded, WeightedBackend{Proxy: proxy, Weight: w})
		}
	}
	return &RoundRobinBalancer{backends: expanded}, nil
}

// ServeHTTP forwards the request to the next backend in rotation.
func (b *RoundRobinBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(b.backends) == 0 {
		http.Error(w, "no backends available", http.StatusServiceUnavailable)
		return
	}
	idx := b.counter.Add(1) - 1
	backend := b.backends[idx%uint64(len(b.backends))]
	backend.Proxy.ServeHTTP(w, r)
}
