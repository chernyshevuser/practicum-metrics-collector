package business

import "context"

type MetricType string

const (
	CounterMT MetricType = "counter"
	GaugeMT   MetricType = "gauge"
)

type Svc interface {
	CollectMetrics(ctx context.Context)
	SendMetrics(ctx context.Context)
	Close()
}
