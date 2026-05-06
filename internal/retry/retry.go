package retry

import (
	"context"
	"time"
)

// Policy defines how retries are performed.
type Policy struct {
	MaxAttempts int
	Delay       time.Duration
	Backoff     float64 // multiplier applied to delay after each attempt
}

// DefaultPolicy returns a sensible default retry policy.
func DefaultPolicy() Policy {
	return Policy{
		MaxAttempts: 3,
		Delay:       500 * time.Millisecond,
		Backoff:     2.0,
	}
}

// Doer is a function that can be retried.
type Doer func(ctx context.Context) error

// Runner executes a Doer with retry logic.
type Runner struct {
	policy Policy
}

// New creates a Runner with the given policy.
func New(p Policy) *Runner {
	return &Runner{policy: p}
}

// Run executes fn up to MaxAttempts times, returning the last error if all
// attempts fail. It respects context cancellation between retries.
func (r *Runner) Run(ctx context.Context, fn Doer) error {
	delay := r.policy.Delay
	var err error
	for attempt := 1; attempt <= r.policy.MaxAttempts; attempt++ {
		if err = fn(ctx); err == nil {
			return nil
		}
		if attempt == r.policy.MaxAttempts {
			break
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
		delay = time.Duration(float64(delay) * r.policy.Backoff)
	}
	return err
}

// Attempts returns the maximum number of attempts configured.
func (r *Runner) Attempts() int {
	return r.policy.MaxAttempts
}
