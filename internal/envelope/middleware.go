package envelope

import (
	"net/http"
	"time"
)

const HeaderEnvelopeService = "X-Envelope-Service"
const HeaderEnvelopeStatus = "X-Envelope-Status"
const HeaderEnvelopeAge = "X-Envelope-Age"

// AttachMiddleware injects envelope metadata into response headers when a
// matching envelope exists for the service named in the request header
// X-Envelope-Service.
func AttachMiddleware(s *Store, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if svc := r.Header.Get(HeaderEnvelopeService); svc != "" {
			if e, ok := s.Get(svc); ok {
				w.Header().Set(HeaderEnvelopeStatus, e.Status)
				age := time.Since(e.WrappedAt).Round(time.Second).String()
				w.Header().Set(HeaderEnvelopeAge, age)
			}
		}
		next.ServeHTTP(w, r)
	})
}
