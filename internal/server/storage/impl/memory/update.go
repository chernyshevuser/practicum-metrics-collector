package memorystorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) UpdateGauge(ctx context.Context, m storage.GaugeMetric) error {
	panic("")
}

func (s *svc) UpdateCounter(ctx context.Context, m storage.CounterMetric) error {
	panic("")
}
