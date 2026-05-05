package aggregator

import (
	"time"

	"github.com/stackprobe/internal/checker"
)

// Status represents the aggregated health status of all services.
type Status struct {
	Timestamp time.Time        `json:"timestamp"`
	OverallOK bool             `json:"overall_ok"`
	Services  []ServiceStatus  `json:"services"`
}

// ServiceStatus holds the result for a single service endpoint.
type ServiceStatus struct {
	Name     string        `json:"name"`
	URL      string        `json:"url"`
	Healthy  bool          `json:"healthy"`
	Latency  time.Duration `json:"latency_ms"`
	Error    string        `json:"error,omitempty"`
}

// Aggregator collects results from the checker and builds a Status report.
type Aggregator struct {
	checker *checker.Checker
}

// New creates a new Aggregator using the provided Checker.
func New(c *checker.Checker) *Aggregator {
	return &Aggregator{checker: c}
}

// Collect runs all health checks and returns an aggregated Status.
func (a *Aggregator) Collect() Status {
	results := a.checker.CheckAll()

	services := make([]ServiceStatus, 0, len(results))
	allOK := true

	for _, r := range results {
		errMsg := ""
		if r.Err != nil {
			errMsg = r.Err.Error()
		}
		if !r.Up {
			allOK = false
		}
		services = append(services, ServiceStatus{
			Name:    r.Name,
			URL:     r.URL,
			Healthy: r.Up,
			Latency: r.Latency,
			Error:   errMsg,
		})
	}

	return Status{
		Timestamp: time.Now().UTC(),
		OverallOK: allOK,
		Services:  services,
	}
}
