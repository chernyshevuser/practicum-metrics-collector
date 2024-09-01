package memorystorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Set(ctx context.Context, metric storage.Metric) (err error) {
	s.storage.Set(storage.BuildKey(metric.ID, metric.Type), metric)
	return nil
}
