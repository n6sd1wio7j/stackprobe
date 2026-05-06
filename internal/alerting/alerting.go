package alerting

import (
	"sync"
	"time"
)

// State tracks the alert state for a single service.
type State int

const (
	StateOK   State = iota
	StateAlert       // service is currently down
)

// Alert represents a fired alert event.
type Alert struct {
	Service   string
	Message   string
	FiredAt   time.Time
	ResolvedAt *time.Time
}

// Manager tracks per-service alert states and fires/resolves alerts.
type Manager struct {
	mu     sync.Mutex
	states map[string]State
	alerts map[string]*Alert
}

// New creates a new alerting Manager.
func New() *Manager {
	return &Manager{
		states: make(map[string]State),
		alerts: make(map[string]*Alert),
	}
}

// Evaluate updates the alert state for a service.
// Returns a fired Alert if state transitions to down, a resolved Alert on recovery, or nil.
func (m *Manager) Evaluate(service string, up bool) *Alert {
	m.mu.Lock()
	defer m.mu.Unlock()

	current := m.states[service]

	if !up && current == StateOK {
		// Transition: OK -> Alert
		m.states[service] = StateAlert
		a := &Alert{
			Service: service,
			Message: service + " is unreachable",
			FiredAt: time.Now().UTC(),
		}
		m.alerts[service] = a
		return a
	}

	if up && current == StateAlert {
		// Transition: Alert -> OK (resolved)
		m.states[service] = StateOK
		a := m.alerts[service]
		if a != nil {
			now := time.Now().UTC()
			a.ResolvedAt = &now
		}
		delete(m.alerts, service)
		return a
	}

	return nil
}

// Active returns all currently firing alerts.
func (m *Manager) Active() []*Alert {
	m.mu.Lock()
	defer m.mu.Unlock()
	alerts := make([]*Alert, 0, len(m.alerts))
	for _, a := range m.alerts {
		alerts = append(alerts, a)
	}
	return alerts
}
