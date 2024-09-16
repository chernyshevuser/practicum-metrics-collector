package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/jackc/pgx/v5"
)

func (s *svc) Get(ctx context.Context, key uint64) (*storage.Metric, error) {
	var metric storage.Metric

	if r := s.wrap(func() error {
		tx, err := s.beginR(ctx)
		if err != nil {
			return fmt.Errorf("failed to create db tx: %w", err)
		}
		defer tx.Rollback(ctx)

		metric, err = s.getQuery(ctx, tx, key)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil
			}
			return err
		}

		if err = tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit db tx: %w", err)
		}

		return nil
	}); r != nil {
		return nil, r
	}

	return &metric, nil
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.Metric, error) {
	var metrics []storage.Metric

	if r := s.wrap(func() error {
		tx, err := s.beginR(ctx)
		if err != nil {
			return fmt.Errorf("failed to create db tx: %w", err)
		}
		defer tx.Rollback(ctx)

		metrics, err = s.getAllQuery(ctx, tx)
		if err != nil {
			return err
		}

		if err = tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit db tx: %w", err)
		}

		return nil
	}); r != nil {
		return nil, r
	}

	return &metrics, nil
}
