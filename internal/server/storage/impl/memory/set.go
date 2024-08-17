package memorystorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Set(ctx context.Context, metrics []storage.Metric) (err error) {
	for _, m := range metrics {
		s.storage.Set(storage.BuildKey(m.ID, m.Type), m)
	}

	return nil
}
