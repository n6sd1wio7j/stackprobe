package checker

import (
	"context"
	"net/http"
	"time"
)

// Status represents the health status of an endpoint.
type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusUnknown Status = "unknown"
)

// Endpoint holds the configuration for a single health-check target.
type Endpoint struct {
	Name string
	URL  string
}

// Result holds the outcome of a single health check.
type Result struct {
	Endpoint   Endpoint
	Status     Status
	StatusCode int
	Latency    time.Duration
	Error      string
	CheckedAt  time.Time
}

// Checker performs HTTP health checks against endpoints.
type Checker struct {
	client  *http.Client
	timeout time.Duration
}

// New creates a new Checker with the given timeout.
func New(timeout time.Duration) *Checker {
	return &Checker{
		client:  &http.Client{Timeout: timeout},
		timeout: timeout,
	}
}

// Check performs an HTTP GET against the endpoint and returns a Result.
func (c *Checker) Check(ctx context.Context, ep Endpoint) Result {
	result := Result{
		Endpoint:  ep,
		Status:    StatusUnknown,
		CheckedAt: time.Now().UTC(),
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, ep.URL, nil)
	if err != nil {
		result.Status = StatusDown
		result.Error = err.Error()
		return result
	}

	start := time.Now()
	resp, err := c.client.Do(req)
	result.Latency = time.Since(start)

	if err != nil {
		result.Status = StatusDown
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Status = StatusUp
	} else {
		result.Status = StatusDown
	}

	return result
}

// CheckAll runs Check concurrently for every endpoint and returns all results.
func (c *Checker) CheckAll(ctx context.Context, endpoints []Endpoint) []Result {
	results := make([]Result, len(endpoints))
	ch := make(chan struct {
		idx int
		res Result
	}, len(endpoints))

	for i, ep := range endpoints {
		go func(idx int, e Endpoint) {
			ch <- struct {
				idx int
				res Result
			}{idx, c.Check(ctx, e)}
		}(i, ep)
	}

	for range endpoints {
		v := <-ch
		results[v.idx] = v.res
	}
	return results
}
