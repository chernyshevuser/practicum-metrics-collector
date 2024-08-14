package business

import (
	"context"
	"fmt"
)

type MetricType string

const (
	Counter MetricType = "counter"
	Gauge   MetricType = "gauge"
	Unknown MetricType = "unknown"
)

var (
	ErrWrongMetricVal = fmt.Errorf("wrong metric value")
)

type MetricsCollector interface {
	UpdateMetric(ctx context.Context, metricType, metricName, metricVal string) error
	UpdateMetricJSON() error
	UpdateMetricsJSON() error
	GetMetricValue() error
	GetMetricValueJSON() error
	GetAllMetrics() error
	PingDB() error
}
