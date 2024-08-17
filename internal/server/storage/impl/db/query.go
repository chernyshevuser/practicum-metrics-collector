package db

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (s *svc) setQuery(ctx context.Context, tx pgx.Tx, metricId, metricType, metricVal string) error {
	q := `
		INSERT INTO Metrics ("id", "type", "val")
		VALUES ($1, $2, $3)
		ON CONFLICT ("id") DO 
		UPDATE SET "val" = $3;
	`

	_, err := tx.Exec(
		ctx,
		q,
		metricId,
		metricType,
		metricVal,
	)

	return err
}

func (s *svc) getQuery(ctx context.Context, tx pgx.Tx, metricId, metricType string) (metricValue string, err error) {
	q := `
		SELECT val
		FROM public.Metrics
		WHERE "id" = $1 AND "type" = $2;
	`
	err = tx.QueryRow(
		ctx,
		q,
		metricId,
		metricType,
	).Scan(&metricValue)
	if err != nil {
		return "", err
	}

	return metricValue, nil
}

type rawMetric struct {
	ID   string
	Type string
	Val  string
}

func (s *svc) getAllQuery(ctx context.Context, tx pgx.Tx) ([]rawMetric, error) {
	q := `
		SELECT *
		FROM public.Metrics;
	`
	rows, err := tx.Query(
		ctx,
		q,
	)
	if err != nil {
		return []rawMetric{}, err
	}
	defer rows.Close()

	var rawMetrics []rawMetric

	for rows.Next() {
		rm := rawMetric{}
		err = rows.Scan(
			&rm.ID,
			&rm.Type,
			&rm.Val,
		)
		if err != nil {
			return []rawMetric{}, err
		}

		rawMetrics = append(rawMetrics, rm)
	}

	err = rows.Err()
	if err != nil {
		return []rawMetric{}, err
	}

	return rawMetrics, nil
}
