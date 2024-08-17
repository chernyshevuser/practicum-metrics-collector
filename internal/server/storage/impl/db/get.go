package db

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Get(ctx context.Context, key string) (*storage.Metric, error) {
	panic("")
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.Metric, error) {
	panic("")
}
