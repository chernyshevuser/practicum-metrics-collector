package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/config"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	memorystorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/memory"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

// TODO add different
func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	fmt.Println("New FROM DB")
	return memorystorage.New(ctx, logger, config.FileStoragePath, config.Restore)
}
