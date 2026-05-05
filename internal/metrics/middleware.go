package metrics

import (
	"net/http"
	"time"

	"github.com/stackprobe/internal/checker"
)

// RecordingChecker wraps a checker and records latency/status metrics.
type RecordingChecker struct {
	inner *checker.Checker
	store *Store
}

// NewRecordingChecker creates a RecordingChecker that records metrics into store.
func NewRecordingChecker(c *checker.Checker, store *Store) *RecordingChecker {
	return &RecordingChecker{inner: c, store: store}
}

// CheckWithMetrics performs a health check and records the result in the store.
func (rc *RecordingChecker) CheckWithMetrics(name, url string) checker.Result {
	start := time.Now()
	result := rc.inner.Check(url)
	latency := time.Since(start)
	rc.store.Record(name, result.Up, latency)
	return result
}

// InstrumentHandler wraps an HTTP handler and records request latency under a given service name.
func InstrumentHandler(name string, store *Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		latency := time.Since(start)
		up := rw.status >= 200 && rw.status < 300
		store.Record(name, up, latency)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}
