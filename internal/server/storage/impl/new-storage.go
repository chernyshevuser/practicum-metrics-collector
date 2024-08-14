package impl

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	memorystorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/memory"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

// TODO add different
func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	return memorystorage.New(ctx, logger)
}
