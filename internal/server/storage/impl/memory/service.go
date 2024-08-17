package memorystorage

import (
	"context"
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	st "github.com/chernyshevuser/practicum-metrics-collector/tools/default-storage"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type svc struct {
	logger   logger.Logger
	storage  *st.Storage
	mu       *sync.Mutex
	filepath string
}

func New(ctx context.Context, logger logger.Logger, filepath string, restoreData bool) (storage.Storage, error) {
	s := svc{
		logger:   logger,
		storage:  st.New[string](),
		mu:       &sync.Mutex{},
		filepath: filepath,
	}

	if restoreData && filepath != "" {
		if err := s.Actualize(ctx); err != nil {
			logger.Errorw(
				"can't actualize memory storage",
				"reason", err,
			)
			return nil, err
		}
		logger.Infow(
			"metrics are actualized successfully",
			"source file", s.filepath,
		)
	}

	return &s, nil
}
