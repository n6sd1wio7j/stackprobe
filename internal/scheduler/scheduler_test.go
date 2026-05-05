package scheduler_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stackprobe/internal/aggregator"
	"github.com/stackprobe/internal/checker"
	"github.com/stackprobe/internal/config"
	"github.com/stackprobe/internal/history"
	"github.com/stackprobe/internal/scheduler"
)

func newHealthyServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func buildConfig(urls []string) *config.Config {
	svcs := make([]config.Service, len(urls))
	for i, u := range urls {
		svcs[i] = config.Service{Name: "svc", URL: u}
	}
	return &config.Config{Services: svcs, Timeout: 2 * time.Second}
}

func TestRun_RecordsHistory(t *testing.T) {
	srv := newHealthyServer()
	defer srv.Close()

	cfg := buildConfig([]string{srv.URL})
	chk := checker.New(cfg)
	agg := aggregator.New(chk, cfg)
	hist := history.New(10)

	sched := scheduler.New(cfg, agg, hist, nil, 50*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	sched.Run(ctx)

	records := hist.Get("svc")
	if len(records) == 0 {
		t.Fatal("expected history records, got none")
	}
	if !records[0].Up {
		t.Errorf("expected service to be up")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	srv := newHealthyServer()
	defer srv.Close()

	cfg := buildConfig([]string{srv.URL})
	chk := checker.New(cfg)
	agg := aggregator.New(chk, cfg)
	hist := history.New(10)

	sched := scheduler.New(cfg, agg, hist, nil, 10*time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		sched.Run(ctx)
		close(done)
	}()

	time.Sleep(30 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancellation")
	}
}
