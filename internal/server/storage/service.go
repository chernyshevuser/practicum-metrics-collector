package storage

import (
	"context"
	"fmt"
)

type Metric struct {
	ID    string
	Type  string
	Val   float64
	Delta int64
}

func BuildKey(metricName, metricType string) string {
	return fmt.Sprintf("%s_%s", metricName, metricType)
}

type Storage interface {
	Set(ctx context.Context, metric Metric) (err error)

	Get(ctx context.Context, key string) (*Metric, error)
	GetAll(ctx context.Context) (*[]Metric, error)

	Lock()
	Unlock()

	Actualize(ctx context.Context) error
	Dump(ctx context.Context) error

	Ping(ctx context.Context) error
	Close() error
}
