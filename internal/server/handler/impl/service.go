package impl

import (
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
)

type api struct {
	mc     business.MetricsCollector
	logger logger.Logger
}

func New(mc business.MetricsCollector, logger logger.Logger) handler.API {
	return &api{
		mc:     mc,
		logger: logger,
	}
}
