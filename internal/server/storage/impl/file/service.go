package filestorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	st "github.com/chernyshevuser/practicum-metrics-collector/tools/storage"
)

type svc struct {
	logger         logger.Logger
	counterStorage *st.Storage
	gaugeStorage   *st.Storage
}

func New(ctx context.Context, logger logger.Logger) (storage.Svc, error) {
	return &svc{
		logger:         logger,
		counterStorage: st.New[string](),
		gaugeStorage:   st.New[string](),
	}, nil
}
