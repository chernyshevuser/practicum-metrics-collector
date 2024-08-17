package impl

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/config"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/db"
	memorystorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/memory"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	if config.DatabaseDsn != "" {
		return db.New(ctx, logger)
	}

	return memorystorage.New(ctx, logger, config.FileStoragePath, config.Restore)
}
