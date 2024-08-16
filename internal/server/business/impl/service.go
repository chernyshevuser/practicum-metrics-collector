package impl

import (
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type svc struct {
	db     storage.Storage
	logger logger.Logger
}

func New(db storage.Storage, logger logger.Logger) business.MetricsCollector {
	return &svc{
		db:     db,
		logger: logger,
	}
}
