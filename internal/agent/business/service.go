package business

import "context"

type MetricType string

const (
	CounterMT MetricType = "counter"
	GaugeMT   MetricType = "gauge"
)

type Agent interface {
	Run(ctx context.Context)
	Close()
}
