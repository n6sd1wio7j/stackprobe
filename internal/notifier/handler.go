package notifier

import (
	"encoding/json"
	"net/http"
)

// WebhookHandler returns an http.Handler that accepts inbound alert
// payloads (e.g. from an external system) and forwards them via Notify.
func (n *Notifier) WebhookHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var p AlertPayload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "invalid JSON body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if p.Service == "" || p.Status == "" {
			http.Error(w, "service and status are required", http.StatusBadRequest)
			return
		}

		if err := n.Notify(p.Service, p.Status, p.Message); err != nil {
			http.Error(w, "failed to forward alert: "+err.Error(), http.StatusBadGateway)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_ = json.NewEncoder(w).Encode(map[string]string{"result": "accepted"})
	})
}
