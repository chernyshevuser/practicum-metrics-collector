package memorystorage

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Get(ctx context.Context, key uint64) (*storage.Metric, error) {
	stored, ok := s.storage.Get(key)
	if !ok {
		return nil, nil
	}

	metric, ok := stored.(storage.Metric)
	if !ok {
		return nil, fmt.Errorf("can't cast stored to storage.Metric")
	}

	return &metric, nil
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.Metric, error) {
	data := s.storage.GetAll()

	res := make([]storage.Metric, 0, len(data))
	for _, tmp := range data {
		metric, ok := tmp.(storage.Metric)
		if !ok {
			return nil, fmt.Errorf("can't cast stored to storage.Metric")
		}

		res = append(res, metric)
	}

	return &res, nil
}
