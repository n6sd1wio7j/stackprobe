package aggregator_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stackprobe/internal/aggregator"
	"github.com/stackprobe/internal/checker"
	"github.com/stackprobe/internal/config"
)

func newServer(t *testing.T, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	}))
}

func TestCollect_AllHealthy(t *testing.T) {
	s1 := newServer(t, http.StatusOK)
	s2 := newServer(t, http.StatusOK)
	defer s1.Close()
	defer s2.Close()

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "svc-a", URL: s1.URL},
			{Name: "svc-b", URL: s2.URL},
		},
	}

	c := checker.New(cfg)
	agg := aggregator.New(c)
	status := agg.Collect()

	if !status.OverallOK {
		t.Error("expected OverallOK to be true when all services are healthy")
	}
	if len(status.Services) != 2 {
		t.Errorf("expected 2 services, got %d", len(status.Services))
	}
	for _, svc := range status.Services {
		if !svc.Healthy {
			t.Errorf("expected service %s to be healthy", svc.Name)
		}
	}
}

func TestCollect_PartialFailure(t *testing.T) {
	s1 := newServer(t, http.StatusOK)
	s2 := newServer(t, http.StatusInternalServerError)
	defer s1.Close()
	defer s2.Close()

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "svc-ok", URL: s1.URL},
			{Name: "svc-fail", URL: s2.URL},
		},
	}

	c := checker.New(cfg)
	agg := aggregator.New(c)
	status := agg.Collect()

	if status.OverallOK {
		t.Error("expected OverallOK to be false when at least one service is unhealthy")
	}

	for _, svc := range status.Services {
		if svc.Name == "svc-fail" && svc.Healthy {
			t.Error("expected svc-fail to be unhealthy")
		}
		if svc.Name == "svc-ok" && !svc.Healthy {
			t.Error("expected svc-ok to be healthy")
		}
	}
}

func TestCollect_TimestampSet(t *testing.T) {
	cfg := &config.Config{
		Services: []config.Service{},
	}
	c := checker.New(cfg)
	agg := aggregator.New(c)
	status := agg.Collect()

	if status.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}
}
