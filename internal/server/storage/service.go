package storage

import (
	"context"

	"github.com/shopspring/decimal"
)

type CounterMetric struct {
	ID    string
	Delta decimal.Decimal
}

type GaugeMetric struct {
	ID    string
	Value decimal.Decimal
}

type Storage interface {
	Update(ctx context.Context, gaugeMetrics []GaugeMetric, counterMetrics []CounterMetric) error

	GetGauge(ctx context.Context, id string) (*GaugeMetric, error)
	GetCounter(ctx context.Context, id string) (*CounterMetric, error)
	GetAll(ctx context.Context) (*[]GaugeMetric, *[]CounterMetric, error)

	Ping(ctx context.Context) error
	Close() error
}
