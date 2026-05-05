package metrics_test

import (
	"testing"
	"time"

	"github.com/stackprobe/internal/metrics"
)

func TestRecord_UpIncrements(t *testing.T) {
	s := metrics.New()
	s.Record("svc-a", true, 50*time.Millisecond)
	s.Record("svc-a", true, 100*time.Millisecond)

	m, ok := s.Get("svc-a")
	if !ok {
		t.Fatal("expected metrics for svc-a")
	}
	if m.TotalChecks != 2 {
		t.Errorf("expected TotalChecks=2, got %d", m.TotalChecks)
	}
	if m.UpCount != 2 {
		t.Errorf("expected UpCount=2, got %d", m.UpCount)
	}
	if m.DownCount != 0 {
		t.Errorf("expected DownCount=0, got %d", m.DownCount)
	}
}

func TestRecord_DownIncrements(t *testing.T) {
	s := metrics.New()
	s.Record("svc-b", false, 200*time.Millisecond)

	m, ok := s.Get("svc-b")
	if !ok {
		t.Fatal("expected metrics for svc-b")
	}
	if m.DownCount != 1 {
		t.Errorf("expected DownCount=1, got %d", m.DownCount)
	}
}

func TestUptimePercent(t *testing.T) {
	s := metrics.New()
	s.Record("svc-c", true, 10*time.Millisecond)
	s.Record("svc-c", true, 10*time.Millisecond)
	s.Record("svc-c", false, 10*time.Millisecond)

	m, _ := s.Get("svc-c")
	got := m.UptimePercent()
	want := 66.66666666666667
	if got < 66.6 || got > 66.7 {
		t.Errorf("expected uptime ~66.7%%, got %.2f", want)
	}
}

func TestAvgLatency(t *testing.T) {
	s := metrics.New()
	s.Record("svc-d", true, 100*time.Millisecond)
	s.Record("svc-d", true, 200*time.Millisecond)

	m, _ := s.Get("svc-d")
	if m.AvgLatency() != 150*time.Millisecond {
		t.Errorf("expected avg latency 150ms, got %v", m.AvgLatency())
	}
}

func TestGet_Unknown(t *testing.T) {
	s := metrics.New()
	_, ok := s.Get("nonexistent")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	s := metrics.New()
	s.Record("alpha", true, 10*time.Millisecond)
	s.Record("beta", false, 20*time.Millisecond)

	all := s.All()
	if len(all) != 2 {
		t.Errorf("expected 2 services, got %d", len(all))
	}
}
