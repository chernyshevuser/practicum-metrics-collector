package db

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Set(ctx context.Context, metric storage.Metric) (err error) {
	tx, err := s.beginW(ctx)
	if err != nil {
		return fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	err = s.setQuery(ctx, tx, storage.BuildKey(metric.ID, metric.Type), metric)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit db tx: %w", err)
	}

	return nil
}
