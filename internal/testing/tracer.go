package testing

import (
	"sync"
	"time"
)

// TraceCollector collects and manages trace events during test execution
type TraceCollector struct {
	mu     sync.RWMutex
	events []TraceEvent
}

// NewTraceCollector creates a new trace collector
func NewTraceCollector() *TraceCollector {
	return &TraceCollector{
		events: make([]TraceEvent, 0),
	}
}

// Start begins tracing for a specific phase/component/action
func (tc *TraceCollector) Start(phase string) *TraceSpan {
	return &TraceSpan{
		phase:     phase,
		component: "unknown",
		action:    "unknown",
		startTime: time.Now(),
		collector: tc,
		details:   make(map[string]interface{}),
	}
}

// StartWithComponent begins tracing with component specified
func (tc *TraceCollector) StartWithComponent(phase, component string) *TraceSpan {
	return &TraceSpan{
		phase:     phase,
		component: component,
		action:    "unknown",
		startTime: time.Now(),
		collector: tc,
		details:   make(map[string]interface{}),
	}
}

// StartWithDetails begins tracing with full details
func (tc *TraceCollector) StartWithDetails(phase, component, action string) *TraceSpan {
	return &TraceSpan{
		phase:     phase,
		component: component,
		action:    action,
		startTime: time.Now(),
		collector: tc,
		details:   make(map[string]interface{}),
	}
}

// AddEvent adds a trace event manually
func (tc *TraceCollector) AddEvent(event TraceEvent) {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.events = append(tc.events, event)
}

// Events returns all collected trace events
func (tc *TraceCollector) Events() []TraceEvent {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	// Return a copy to prevent external modification
	events := make([]TraceEvent, len(tc.events))
	copy(events, tc.events)
	return events
}

// Clear clears all trace events
func (tc *TraceCollector) Clear() {
	tc.mu.Lock()
	defer tc.mu.Unlock()
	tc.events = tc.events[:0]
}

// Count returns the number of trace events
func (tc *TraceCollector) Count() int {
	tc.mu.RLock()
	defer tc.mu.RUnlock()
	return len(tc.events)
}

// FilterByPhase returns events filtered by phase
func (tc *TraceCollector) FilterByPhase(phase string) []TraceEvent {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	var filtered []TraceEvent
	for _, event := range tc.events {
		if event.Phase == phase {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// FilterByComponent returns events filtered by component
func (tc *TraceCollector) FilterByComponent(component string) []TraceEvent {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	var filtered []TraceEvent
	for _, event := range tc.events {
		if event.Component == component {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// FilterByStatus returns events filtered by status
func (tc *TraceCollector) FilterByStatus(status string) []TraceEvent {
	tc.mu.RLock()
	defer tc.mu.RUnlock()

	var filtered []TraceEvent
	for _, event := range tc.events {
		if event.Status == status {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

// TraceSpan represents an active trace span
type TraceSpan struct {
	phase     string
	component string
	action    string
	startTime time.Time
	collector *TraceCollector
	details   map[string]interface{}
	finished  bool
}

// WithComponent sets the component for this span
func (ts *TraceSpan) WithComponent(component string) *TraceSpan {
	ts.component = component
	return ts
}

// WithAction sets the action for this span
func (ts *TraceSpan) WithAction(action string) *TraceSpan {
	ts.action = action
	return ts
}

// WithDetails adds details to this span
func (ts *TraceSpan) WithDetails(details map[string]interface{}) *TraceSpan {
	for k, v := range details {
		ts.details[k] = v
	}
	return ts
}

// WithDetail adds a single detail to this span
func (ts *TraceSpan) WithDetail(key string, value interface{}) *TraceSpan {
	ts.details[key] = value
	return ts
}

// End completes this trace span with success status
func (ts *TraceSpan) End() {
	ts.EndWithStatus("success", nil)
}

// EndWithError completes this trace span with error status
func (ts *TraceSpan) EndWithError(err error) {
	ts.EndWithStatus("error", err)
}

// EndWithStatus completes this trace span with specified status
func (ts *TraceSpan) EndWithStatus(status string, err error) {
	if ts.finished {
		return // Already finished
	}

	ts.finished = true
	duration := time.Since(ts.startTime)

	event := TraceEvent{
		Timestamp: ts.startTime,
		Phase:     ts.phase,
		Component: ts.component,
		Action:    ts.action,
		Duration:  duration,
		Status:    status,
		Details:   ts.details,
	}

	if err != nil {
		event.Error = err.Error()
	}

	ts.collector.AddEvent(event)
}

// IsFinished returns whether this span is finished
func (ts *TraceSpan) IsFinished() bool {
	return ts.finished
}

// Duration returns the duration of this span (if finished)
func (ts *TraceSpan) Duration() time.Duration {
	if !ts.finished {
		return time.Since(ts.startTime)
	}
	return 0
}

// Phase returns the phase of this span
func (ts *TraceSpan) Phase() string {
	return ts.phase
}

// Component returns the component of this span
func (ts *TraceSpan) Component() string {
	return ts.component
}

// Action returns the action of this span
func (ts *TraceSpan) Action() string {
	return ts.action
}

// Details returns the details of this span
func (ts *TraceSpan) Details() map[string]interface{} {
	// Return a copy to prevent external modification
	details := make(map[string]interface{})
	for k, v := range ts.details {
		details[k] = v
	}
	return details
}
