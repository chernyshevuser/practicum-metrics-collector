package storage

import (
	"context"
)

type Metric struct {
	ID    string
	Type  string
	Val   float64
	Delta int64
}

func BuildKey(metricName, metricType string) uint64 {
	var hash uint64 = 0
	const (
		prime uint64 = 31
		mod   uint64 = 1e9 + 7
	)

	for i := 0; i < len(metricName); i++ {
		hash = (hash*prime + uint64(metricName[i])) % mod
	}
	for i := 0; i < len(metricType); i++ {
		hash = (hash*prime + uint64(metricType[i])) % mod
	}

	return hash
}

type Storage interface {
	Set(ctx context.Context, metric Metric) (err error)

	Get(ctx context.Context, key uint64) (*Metric, error)
	GetAll(ctx context.Context) (*[]Metric, error)

	Lock()
	Unlock()

	Actualize(ctx context.Context) error
	Dump(ctx context.Context) error

	Ping(ctx context.Context) error
	Close() error
}
