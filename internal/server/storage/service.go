package storage

import "context"

type CounterMetric struct {
	ID    string
	Delta int64
}

type GaugeMetric struct {
	ID    string
	Value float64
}

type Svc interface {
	UpdateGauge(ctx context.Context, m GaugeMetric) error
	UpdateCounter(ctx context.Context, m CounterMetric) error

	GetGauge(ctx context.Context, id string) (GaugeMetric, error)
	GetCounter(ctx context.Context, id string) (CounterMetric, error)
	GetAll(ctx context.Context) ([]GaugeMetric, []CounterMetric, error)

	Ping(ctx context.Context) error
	Actualize(ctx context.Context) error

	Close() error
}
