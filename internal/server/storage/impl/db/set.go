package db

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) Set(ctx context.Context, metrics []storage.Metric) (err error) {
	tx, err := s.beginW(ctx)
	if err != nil {
		return fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	for _, m := range metrics {
		err := s.setQuery(ctx, tx, m.ID, m.Type, m.Val.String())
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit db tx: %w", err)
	}

	return nil
}
