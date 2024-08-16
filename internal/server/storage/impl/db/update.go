package db

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Update(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) (err error) {
	panic("")
}
