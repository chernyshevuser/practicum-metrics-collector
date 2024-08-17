package storage

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type Metric struct {
	ID   string
	Type string
	Val  decimal.Decimal
}

func BuildKey(metricId, metricType string) string {
	return fmt.Sprintf("%s_%s", metricId, metricType)
}

type Storage interface {
	Set(ctx context.Context, metrics []Metric) (err error)

	Get(ctx context.Context, key string) (*Metric, error)
	GetAll(ctx context.Context) (*[]Metric, error)

	Lock()
	Unlock()

	Actualize(ctx context.Context) error
	Dump(ctx context.Context) error

	Ping(ctx context.Context) error
	Close() error
}
