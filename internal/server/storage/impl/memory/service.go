package memorystorage

import (
	"context"
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	st "github.com/chernyshevuser/practicum-metrics-collector/tools/default-storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type svc struct {
	logger         logger.Logger
	counterStorage *st.Storage
	gaugeStorage   *st.Storage
	mu             *sync.Mutex
}

func New(ctx context.Context, logger logger.Logger) (storage.Storage, error) {
	return &svc{
		logger:         logger,
		counterStorage: st.New[string](),
		gaugeStorage:   st.New[string](),
		mu:             &sync.Mutex{},
	}, nil
}
