package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// AlertPayload is the JSON body sent to a webhook endpoint.
type AlertPayload struct {
	Service   string    `json:"service"`
	Status    string    `json:"status"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

// Notifier sends alert notifications to a configured webhook URL.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a Notifier that posts alerts to webhookURL.
func New(webhookURL string, timeout time.Duration) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: timeout},
	}
}

// Notify sends an alert payload for the given service to the webhook.
func (n *Notifier) Notify(service, status, message string) error {
	payload := AlertPayload{
		Service:   service,
		Status:    status,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notifier: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: post to webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		return fmt.Errorf("notifier: webhook returned non-2xx status %d: %s", resp.StatusCode, bytes.TrimSpace(respBody))
	}

	return nil
}
