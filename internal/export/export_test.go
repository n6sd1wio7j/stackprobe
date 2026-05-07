package export_test

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/stackprobe/internal/export"
)

type stubSource struct{ recs []export.Record }

func (s *stubSource) Records() []export.Record { return s.recs }

func sampleRecords() []export.Record {
	return []export.Record{
		{Service: "api", Status: "up", LatencyMs: 42, CheckedAt: time.Unix(0, 0).UTC()},
		{Service: "db", Status: "down", LatencyMs: 0, CheckedAt: time.Unix(0, 0).UTC()},
	}
}

func TestWriteJSON(t *testing.T) {
	exp := export.New(&stubSource{recs: sampleRecords()})
	var buf bytes.Buffer
	if err := exp.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON: %v", err)
	}
	var out []export.Record
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("want 2 records, got %d", len(out))
	}
	if out[0].Service != "api" || out[1].Status != "down" {
		t.Errorf("unexpected records: %+v", out)
	}
}

func TestWriteCSV(t *testing.T) {
	exp := export.New(&stubSource{recs: sampleRecords()})
	var buf bytes.Buffer
	if err := exp.WriteCSV(&buf); err != nil {
		t.Fatalf("WriteCSV: %v", err)
	}
	r := csv.NewReader(&buf)
	rows, err := r.ReadAll()
	if err != nil {
		t.Fatalf("csv read: %v", err)
	}
	// header + 2 data rows
	if len(rows) != 3 {
		t.Fatalf("want 3 rows, got %d", len(rows))
	}
	if rows[0][0] != "service" {
		t.Errorf("expected header 'service', got %q", rows[0][0])
	}
	if rows[1][0] != "api" || rows[2][1] != "down" {
		t.Errorf("unexpected data rows: %v", rows[1:])
	}
}

func TestHandler_JSONDefault(t *testing.T) {
	exp := export.New(&stubSource{recs: sampleRecords()})
	h := export.NewHandler(exp)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/export", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "application/json") {
		t.Errorf("expected JSON content-type")
	}
}

func TestHandler_CSVFormat(t *testing.T) {
	exp := export.New(&stubSource{recs: sampleRecords()})
	h := export.NewHandler(exp)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/export?format=csv", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d", rec.Code)
	}
	if !strings.Contains(rec.Header().Get("Content-Type"), "text/csv") {
		t.Errorf("expected CSV content-type")
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	exp := export.New(&stubSource{})
	h := export.NewHandler(exp)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/export", nil)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("want 405, got %d", rec.Code)
	}
}
