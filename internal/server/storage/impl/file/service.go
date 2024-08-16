package filestorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	st "github.com/chernyshevuser/practicum-metrics-collector/tools/default-storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type svc struct {
	logger         logger.Logger
	counterStorage *st.Storage
	gaugeStorage   *st.Storage
}

func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	return &svc{
		logger:         logger,
		counterStorage: st.New[string](),
		gaugeStorage:   st.New[string](),
	}, nil
}
