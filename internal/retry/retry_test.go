package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stackprobe/internal/retry"
)

var errFail = errors.New("fail")

func TestRun_SucceedsOnFirstAttempt(t *testing.T) {
	r := retry.New(retry.DefaultPolicy())
	calls := 0
	err := r.Run(context.Background(), func(_ context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestRun_RetriesUntilSuccess(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Backoff: 1.0}
	r := retry.New(p)
	var calls int32
	err := r.Run(context.Background(), func(_ context.Context) error {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			return errFail
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_ReturnsLastError(t *testing.T) {
	p := retry.Policy{MaxAttempts: 3, Delay: time.Millisecond, Backoff: 1.0}
	r := retry.New(p)
	var calls int32
	err := r.Run(context.Background(), func(_ context.Context) error {
		atomic.AddInt32(&calls, 1)
		return errFail
	})
	if !errors.Is(err, errFail) {
		t.Fatalf("expected errFail, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

func TestRun_RespectsContextCancellation(t *testing.T) {
	p := retry.Policy{MaxAttempts: 5, Delay: 100 * time.Millisecond, Backoff: 1.0}
	r := retry.New(p)
	ctx, cancel := context.WithCancel(context.Background())
	var calls int32
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := r.Run(ctx, func(_ context.Context) error {
		atomic.AddInt32(&calls, 1)
		return errFail
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls >= 5 {
		t.Fatal("context cancellation did not stop retries")
	}
}

func TestAttempts(t *testing.T) {
	p := retry.Policy{MaxAttempts: 7, Delay: time.Millisecond, Backoff: 1.0}
	r := retry.New(p)
	if r.Attempts() != 7 {
		t.Fatalf("expected 7, got %d", r.Attempts())
	}
}
