package business

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
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

type CounterMetric struct {
	ID    string
	Delta decimal.Decimal
}

type GaugeMetric struct {
	ID    string
	Value decimal.Decimal
}

type MetricsCollector interface {
	UpdateMetric(ctx context.Context, metricType, metricName, metricVal string) error
	UpdateMetricJSON() error
	UpdateMetricsJSON() error
	GetMetricValue(ctx context.Context, metricType, metricName string) (*decimal.Decimal, error)
	GetMetricValueJSON() error
	GetAllMetrics(ctx context.Context) ([]CounterMetric, []GaugeMetric, error)
	PingDB() error
	Close()
}
