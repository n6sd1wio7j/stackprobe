package scheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/stackprobe/internal/aggregator"
	"github.com/stackprobe/internal/config"
	"github.com/stackprobe/internal/history"
	"github.com/stackprobe/internal/notifier"
)

// Scheduler periodically runs health checks and records results.
type Scheduler struct {
	cfg        *config.Config
	agg        *aggregator.Aggregator
	hist       *history.History
	notif      *notifier.Notifier
	interval   time.Duration
	prevStatus map[string]bool
	mu         sync.Mutex
}

// New creates a new Scheduler.
func New(cfg *config.Config, agg *aggregator.Aggregator, hist *history.History, notif *notifier.Notifier, interval time.Duration) *Scheduler {
	return &Scheduler{
		cfg:        cfg,
		agg:        agg,
		hist:       hist,
		notif:      notif,
		interval:   interval,
		prevStatus: make(map[string]bool),
	}
}

// Run starts the scheduler loop, blocking until ctx is cancelled.
func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Run once immediately before waiting for the first tick.
	s.runOnce(ctx)

	for {
		select {
		case <-ticker.C:
			s.runOnce(ctx)
		case <-ctx.Done():
			log.Println("scheduler: stopping")
			return
		}
	}
}

func (s *Scheduler) runOnce(ctx context.Context) {
	results := s.agg.Collect(ctx)

	s.mu.Lock()
	defer s.mu.Unlock()

	for _, r := range results {
		s.hist.Add(r)

		prev, seen := s.prevStatus[r.Name]
		if seen && prev != r.Up {
			if s.notif != nil {
				if err := s.notif.Notify(ctx, r); err != nil {
					log.Printf("scheduler: notify error for %s: %v", r.Name, err)
				}
			}
		}
		s.prevStatus[r.Name] = r.Up
	}
}
