package capacity

import (
	"net/http"
	"strconv"
)

// SaturationHeader is the response header set by GuardMiddleware.
const SaturationHeader = "X-Capacity-Saturated"

// GuardMiddleware rejects requests with 503 when the named service is
// saturated according to the Store. The service name is resolved via
// the provided resolver function so callers can derive it from the
// request (e.g. from a path segment or header).
func GuardMiddleware(s *Store, resolve func(*http.Request) string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		svc := resolve(r)
		if svc != "" {
			rec, ok := s.Get(svc)
			if ok {
				w.Header().Set(SaturationHeader, strconv.FormatBool(s.IsSaturated(svc)))
				w.Header().Set("X-Capacity-Percent", strconv.FormatFloat(rec.Percent, 'f', 1, 64))
				if s.IsSaturated(svc) {
					http.Error(w, "service at capacity", http.StatusServiceUnavailable)
					return
				}
			}
		}
		next.ServeHTTP(w, r)
	})
}
