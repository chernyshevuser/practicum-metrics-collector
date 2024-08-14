package impl

import (
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

type svc struct {
	db storage.Storage
}

func New(db storage.Storage) business.MetricsCollector {
	return &svc{
		db: db,
	}
}
