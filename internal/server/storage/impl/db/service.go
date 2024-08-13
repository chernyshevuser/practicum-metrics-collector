package db

import (
	"context"

	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/config"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5/pgxpool"
)

type svc struct {
	conn   *pgxpool.Pool
	logger logger.Logger
}

func New(ctx context.Context, logger logger.Logger) (svc storage.Svc, err error) {
	dbPool, err := pgxpool.New(
		ctx,
		config.DatabaseURI,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create dbPool: %w", err)
	}

	err = dbPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgreSQL: %w", err)
	}

	panic("")

	// svc = &svc{
	// 	conn:   dbPool,
	// 	logger: logger,
	// }

	// return svc, nil
}
