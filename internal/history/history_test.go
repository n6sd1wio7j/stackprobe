package history_test

import (
	"testing"
	"time"

	"github.com/stackprobe/internal/history"
)

func makeRecord(service string, healthy bool) history.Record {
	return history.Record{
		Service:   service,
		Healthy:   healthy,
		Status:    200,
		CheckedAt: time.Now(),
	}
}

func TestAdd_And_Get(t *testing.T) {
	s := history.New(5)
	s.Add(makeRecord("api", true))
	s.Add(makeRecord("api", false))

	recs := s.Get("api")
	if len(recs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(recs))
	}
	if recs[0].Healthy != true {
		t.Errorf("expected first record healthy=true")
	}
	if recs[1].Healthy != false {
		t.Errorf("expected second record healthy=false")
	}
}

func TestAdd_EnforcesLimit(t *testing.T) {
	const limit = 3
	s := history.New(limit)

	for i := 0; i < 10; i++ {
		s.Add(makeRecord("svc", i%2 == 0))
	}

	recs := s.Get("svc")
	if len(recs) != limit {
		t.Fatalf("expected %d records after eviction, got %d", limit, len(recs))
	}
}

func TestGet_UnknownService(t *testing.T) {
	s := history.New(5)
	recs := s.Get("nonexistent")
	if recs == nil {
		t.Fatal("expected non-nil slice for unknown service")
	}
	if len(recs) != 0 {
		t.Fatalf("expected empty slice, got %d records", len(recs))
	}
}

func TestAll_MultipleServices(t *testing.T) {
	s := history.New(5)
	s.Add(makeRecord("alpha", true))
	s.Add(makeRecord("beta", false))
	s.Add(makeRecord("alpha", false))

	all := s.All()
	if len(all["alpha"]) != 2 {
		t.Errorf("expected 2 records for alpha, got %d", len(all["alpha"]))
	}
	if len(all["beta"]) != 1 {
		t.Errorf("expected 1 record for beta, got %d", len(all["beta"]))
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	s := history.New(5)
	s.Add(makeRecord("svc", true))

	all := s.All()
	all["svc"][0].Healthy = false

	original := s.Get("svc")
	if !original[0].Healthy {
		t.Error("All() should return a copy; mutating it must not affect the store")
	}
}
