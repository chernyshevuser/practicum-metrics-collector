package impl

import (
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/constants"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"

	"github.com/shopspring/decimal"
)

type Metric struct {
	ID   string
	Type string
	Val  decimal.Decimal
}

type svc struct {
	logger logger.Logger

	metrics   []Metric
	pollCount decimal.Decimal
	mu        *sync.Mutex

	ch chan Metric
	wg *sync.WaitGroup
}

func New(logger logger.Logger) business.Svc {

	return &svc{
		logger: logger,

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		ch: make(chan Metric, constants.Sz),
		wg: &sync.WaitGroup{},
	}
}
