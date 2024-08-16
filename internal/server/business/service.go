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

type RawMetric struct {
	ID    string
	Type  string
	Value string
}

type MetricsCollector interface {
	UpdateMetrics(ctx context.Context, metrics []RawMetric) (updatedCounterMetric []CounterMetric, updatedGaugeMetrics []GaugeMetric, err error)
	GetMetricValue(ctx context.Context, metricType, metricName string) (val *decimal.Decimal, mType MetricType, err error)
	GetAllMetrics(ctx context.Context) ([]CounterMetric, []GaugeMetric, error)
	PingDB() error
	Close()
}
