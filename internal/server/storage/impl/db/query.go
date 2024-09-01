package db

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/jackc/pgx/v5"
)

func (s *svc) setQuery(ctx context.Context, tx pgx.Tx, key string, metric storage.Metric) error {
	q := `
		INSERT INTO Metrics ("key", "id", "type", "val", "delta")
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT ("key") DO 
		UPDATE SET "val" = $4, "delta" = $5;
	`

	_, err := tx.Exec(
		ctx,
		q,
		key,
		metric.ID,
		metric.Type,
		metric.Val,
		metric.Delta,
	)

	return err
}

func (s *svc) getQuery(ctx context.Context, tx pgx.Tx, key string) (metric storage.Metric, err error) {
	q := `
		SELECT "id", "type", "val", "delta"
		FROM public.Metrics
		WHERE "key" = $1;
	`

	err = tx.QueryRow(
		ctx,
		q,
		key,
	).Scan(
		&metric.Type,
		&metric.ID,
		&metric.Val,
		&metric.Delta,
	)
	if err != nil {
		return storage.Metric{}, err
	}

	return metric, nil
}

func (s *svc) getAllQuery(ctx context.Context, tx pgx.Tx) ([]storage.Metric, error) {
	q := `
		SELECT "id", "type", "val", "delta"
		FROM public.Metrics;
	`
	rows, err := tx.Query(
		ctx,
		q,
	)
	if err != nil {
		return []storage.Metric{}, err
	}
	defer rows.Close()

	var metrics []storage.Metric

	for rows.Next() {
		m := storage.Metric{}
		err = rows.Scan(
			&m.ID,
			&m.Type,
			&m.Val,
			&m.Delta,
		)
		if err != nil {
			return []storage.Metric{}, err
		}

		metrics = append(metrics, m)
	}

	err = rows.Err()
	if err != nil {
		return []storage.Metric{}, err
	}

	return metrics, nil
}
