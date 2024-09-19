// Package business defines business logic related to agents that manage metrics processing.
// It includes metric types and an Agent interface for running and closing agents.
package business

import "context"

// MetricType defines a string type for identifying different types of metrics, such as Counter and Gauge.
type MetricType string

const (
	// CounterMT represents a metric type for counters, which are cumulative metrics.
	CounterMT MetricType = "counter"
	// GaugeMT represents a metric type for gauges, which store real-time values.
	GaugeMT MetricType = "gauge"
)

// Agent defines an interface for managing agents that handle metrics processing.
type Agent interface {
	// Run starts the agent's operation in the given context.
	// The agent can perform continuous or scheduled tasks as long as the context is active.
	Run(ctx context.Context)

	// Close gracefully shuts down the agent, releasing any resources and stopping its operations.
	Close()
}
