package export

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// Record represents a single exportable health-check snapshot.
type Record struct {
	Service   string        `json:"service"`
	Status    string        `json:"status"`
	LatencyMs int64         `json:"latency_ms"`
	CheckedAt time.Time     `json:"checked_at"`
}

// Source is anything that can supply a slice of Records.
type Source interface {
	Records() []Record
}

// Exporter writes Records in a requested format.
type Exporter struct {
	src Source
}

// New returns an Exporter backed by src.
func New(src Source) *Exporter {
	return &Exporter{src: src}
}

// WriteJSON serialises all records as a JSON array to w.
func (e *Exporter) WriteJSON(w io.Writer) error {
	recs := e.src.Records()
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(recs); err != nil {
		return fmt.Errorf("export: json encode: %w", err)
	}
	return nil
}

// WriteCSV serialises all records as CSV to w.
// Header row: service,status,latency_ms,checked_at
func (e *Exporter) WriteCSV(w io.Writer) error {
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"service", "status", "latency_ms", "checked_at"}); err != nil {
		return fmt.Errorf("export: csv header: %w", err)
	}
	for _, r := range e.src.Records() {
		row := []string{
			r.Service,
			r.Status,
			fmt.Sprintf("%d", r.LatencyMs),
			r.CheckedAt.UTC().Format(time.RFC3339),
		}
		if err := cw.Write(row); err != nil {
			return fmt.Errorf("export: csv row: %w", err)
		}
	}
	cw.Flush()
	return cw.Error()
}
