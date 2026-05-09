package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"stackprobe/internal/snapshot"
)

func makeRecord(svc string, healthy bool) snapshot.Record {
	return snapshot.Record{
		Service:    svc,
		Healthy:    healthy,
		Latency:    42 * time.Millisecond,
		CapturedAt: time.Now().UTC(),
	}
}

func TestSave_And_Get(t *testing.T) {
	s := snapshot.New()
	r := makeRecord("api", true)
	s.Save(r)

	got, ok := s.Get("api")
	if !ok {
		t.Fatal("expected record to exist")
	}
	if got.Service != "api" || !got.Healthy {
		t.Fatalf("unexpected record: %+v", got)
	}
}

func TestGet_Unknown(t *testing.T) {
	s := snapshot.New()
	_, ok := s.Get("unknown")
	if ok {
		t.Fatal("expected no record for unknown service")
	}
}

func TestSave_Overwrites(t *testing.T) {
	s := snapshot.New()
	s.Save(makeRecord("db", true))
	s.Save(makeRecord("db", false))

	got, _ := s.Get("db")
	if got.Healthy {
		t.Fatal("expected overwritten record to be unhealthy")
	}
}

func TestAll_ReturnsAllServices(t *testing.T) {
	s := snapshot.New()
	s.Save(makeRecord("svc-a", true))
	s.Save(makeRecord("svc-b", false))

	all := s.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 records, got %d", len(all))
	}
}

func TestSave_SetsTimestampWhenZero(t *testing.T) {
	s := snapshot.New()
	r := snapshot.Record{Service: "x", Healthy: true}
	s.Save(r)

	got, _ := s.Get("x")
	if got.CapturedAt.IsZero() {
		t.Fatal("expected CapturedAt to be set automatically")
	}
}

func TestDump_WritesJSON(t *testing.T) {
	s := snapshot.New()
	s.Save(makeRecord("alpha", true))

	tmp := filepath.Join(t.TempDir(), "snap.json")
	if err := s.Dump(tmp); err != nil {
		t.Fatalf("Dump error: %v", err)
	}

	data, err := os.ReadFile(tmp)
	if err != nil {
		t.Fatalf("ReadFile error: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected non-empty JSON file")
	}
}
