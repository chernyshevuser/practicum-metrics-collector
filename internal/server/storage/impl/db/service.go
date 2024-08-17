package db

import (
	"context"

	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/config"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type svc struct {
	conn   *pgxpool.Pool
	logger logger.Logger
}

func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	dbPool, err := pgxpool.New(
		ctx,
		config.DatabaseDsn,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create dbPool: %v", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgreSQL: %v", err)
	}

	s := svc{
		conn:   dbPool,
		logger: logger,
	}

	err = s.Actualize(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't actualize db, reason: %v", err)
	}

	return &s, nil
}
