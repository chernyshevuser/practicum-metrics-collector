// Package business defines business logic for handling and processing metrics,
// including counters, gauges, and raw metric types.
package business

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// MetricType defines a type for metric categories such as Counter, Gauge, and Unknown.
type MetricType string

const (
	// Counter represents a metric type for counters, which are cumulative values.
	Counter MetricType = "counter"
	// Gauge represents a metric type for gauges, which store real-time values.
	Gauge MetricType = "gauge"
	// Unknown represents an unspecified or unrecognized metric type.
	Unknown MetricType = "unknown"
)

var (
	// ErrWrongMetricVal is returned when a metric value is invalid or cannot be processed.
	ErrWrongMetricVal = fmt.Errorf("wrong metric value")
)

// CounterMetric represents a counter metric with an ID and a cumulative Delta value.
type CounterMetric struct {
	ID    string          // ID is the identifier of the counter metric.
	Delta decimal.Decimal // Delta is the accumulated value of the counter metric.
}

// GaugeMetric represents a gauge metric with an ID and a real-time Value.
type GaugeMetric struct {
	ID    string          // ID is the identifier of the gauge metric.
	Value decimal.Decimal // Value is the current real-time value of the gauge metric.
}

// RawMetric represents a raw metric received, typically in string form, before being processed.
type RawMetric struct {
	ID    string // ID is the identifier of the raw metric.
	Type  string // Type indicates the metric type (e.g., counter or gauge).
	Value string // Value is the raw string value of the metric.
}

// MetricsCollector defines the interface for managing metrics, including updating and retrieving them.
type MetricsCollector interface {
	// UpdateMetrics processes a list of raw metrics and updates the internal state.
	UpdateMetrics(ctx context.Context, metrics []RawMetric) (updatedCounterMetric []CounterMetric, updatedGaugeMetrics []GaugeMetric, err error)

	// GetMetricValue retrieves the value of a specific metric given its type and name.
	GetMetricValue(ctx context.Context, metricType, metricName string) (val *decimal.Decimal, mType MetricType, err error)

	// GetAllMetrics returns all counter and gauge metrics currently stored in the system.
	GetAllMetrics(ctx context.Context) ([]CounterMetric, []GaugeMetric, error)

	// PingDB checks the connectivity and responsiveness of the database or storage backend.
	PingDB(ctx context.Context) error

	// Close releases any resources held by the MetricsCollector and performs cleanup operations.
	Close()
}
