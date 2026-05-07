package maintenance_test

import (
	"testing"
	"time"

	"stackprobe/internal/maintenance"
)

func makeWindow(service string, start, end time.Time) maintenance.Window {
	return maintenance.Window{
		Service:  service,
		Reason:   "planned downtime",
		StartsAt: start,
		EndsAt:   end,
	}
}

func TestIsActive_DuringWindow(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-a", now.Add(-time.Minute), now.Add(time.Minute)))

	if !s.IsActive("svc-a", now) {
		t.Fatal("expected service to be in maintenance")
	}
}

func TestIsActive_BeforeWindow(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-b", now.Add(time.Hour), now.Add(2*time.Hour)))

	if s.IsActive("svc-b", now) {
		t.Fatal("window has not started yet; should not be active")
	}
}

func TestIsActive_AfterWindow(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-c", now.Add(-2*time.Hour), now.Add(-time.Hour)))

	if s.IsActive("svc-c", now) {
		t.Fatal("window has ended; should not be active")
	}
}

func TestIsActive_UnknownService(t *testing.T) {
	s := maintenance.New()
	if s.IsActive("unknown", time.Now()) {
		t.Fatal("unknown service must never be active")
	}
}

func TestDelete_RemovesWindow(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-d", now.Add(-time.Minute), now.Add(time.Minute)))
	s.Delete("svc-d")

	if s.IsActive("svc-d", now) {
		t.Fatal("window should have been removed")
	}
}

func TestAll_ReturnsAllWindows(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-e", now, now.Add(time.Hour)))
	s.Set(makeWindow("svc-f", now, now.Add(time.Hour)))

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 windows, got %d", len(all))
	}
}

func TestSet_Overwrites(t *testing.T) {
	s := maintenance.New()
	now := time.Now()
	s.Set(makeWindow("svc-g", now.Add(-time.Hour), now.Add(-time.Minute)))
	// Overwrite with an active window.
	s.Set(makeWindow("svc-g", now.Add(-time.Minute), now.Add(time.Hour)))

	if !s.IsActive("svc-g", now) {
		t.Fatal("overwritten window should now be active")
	}
}
