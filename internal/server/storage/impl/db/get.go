package db

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) Get(ctx context.Context, key string) (*storage.Metric, error) {
	tx, err := s.beginR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	metricId, metricType, err := storage.ParseKey(key)
	if err != nil {
		return nil, err
	}

	rawValue, err := s.getQuery(ctx, tx, metricId, metricType)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit db tx: %w", err)
	}

	val, err := decimal.NewFromString(rawValue)
	if err != nil {
		return nil, err
	}

	return &storage.Metric{
		ID:   metricId,
		Type: metricType,
		Val:  val,
	}, nil
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.Metric, error) {
	tx, err := s.beginR(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create db tx: %w", err)
	}
	defer tx.Rollback(ctx)

	rawMetrics, err := s.getAllQuery(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit db tx: %w", err)
	}

	var metrics []storage.Metric

	for _, m := range rawMetrics {
		val, err := decimal.NewFromString(m.Val)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, storage.Metric{
			ID:   m.ID,
			Type: m.Type,
			Val:  val,
		})
	}

	return &metrics, nil
}
