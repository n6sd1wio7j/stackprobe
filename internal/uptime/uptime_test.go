package uptime

import (
	"testing"
	"time"
)

func TestReport_InitialUp(t *testing.T) {
	tr := New()
	tr.Report("svc-a", true)
	r, ok := tr.Get("svc-a")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if !r.IsUp {
		t.Error("expected service to be up")
	}
	if r.DownSince != nil {
		t.Error("expected DownSince to be nil when up")
	}
}

func TestReport_InitialDown(t *testing.T) {
	tr := New()
	tr.Report("svc-b", false)
	r, ok := tr.Get("svc-b")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if r.IsUp {
		t.Error("expected service to be down")
	}
	if r.DownSince == nil {
		t.Error("expected DownSince to be set")
	}
}

func TestReport_TransitionUpToDown(t *testing.T) {
	tr := New()
	tr.Report("svc-c", true)
	time.Sleep(2 * time.Millisecond)
	tr.Report("svc-c", false)
	r, _ := tr.Get("svc-c")
	if r.IsUp {
		t.Error("expected service to be down after transition")
	}
	if r.DownSince == nil {
		t.Error("expected DownSince to be set after transition")
	}
}

func TestReport_TransitionDownToUp(t *testing.T) {
	tr := New()
	tr.Report("svc-d", false)
	time.Sleep(2 * time.Millisecond)
	tr.Report("svc-d", true)
	r, _ := tr.Get("svc-d")
	if !r.IsUp {
		t.Error("expected service to be up after recovery")
	}
	if r.DownSince != nil {
		t.Error("expected DownSince to be nil after recovery")
	}
}

func TestGet_Unknown(t *testing.T) {
	tr := New()
	_, ok := tr.Get("unknown")
	if ok {
		t.Error("expected false for unknown service")
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	tr := New()
	tr.Report("svc-1", true)
	tr.Report("svc-2", false)
	all := tr.All()
	if len(all) != 2 {
		t.Errorf("expected 2 records, got %d", len(all))
	}
}

func TestDuration_Up(t *testing.T) {
	tr := New()
	tr.Report("svc-e", true)
	time.Sleep(5 * time.Millisecond)
	d, ok := tr.Duration("svc-e")
	if !ok {
		t.Fatal("expected duration to be available")
	}
	if d < 5*time.Millisecond {
		t.Errorf("expected duration >= 5ms, got %v", d)
	}
}

func TestDuration_Unknown(t *testing.T) {
	tr := New()
	_, ok := tr.Duration("ghost")
	if ok {
		t.Error("expected false for unknown service")
	}
}
