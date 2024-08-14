package db

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) GetGauge(ctx context.Context, id string) (*storage.GaugeMetric, error) {
	panic("")
}

func (s *svc) GetCounter(ctx context.Context, id string) (*storage.CounterMetric, error) {
	panic("")
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.GaugeMetric, *[]storage.CounterMetric, error) {
	panic("")
}
