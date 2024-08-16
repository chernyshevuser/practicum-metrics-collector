package impl

import (
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business"
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

	closeCh chan struct{}
	wg      *sync.WaitGroup

	updateInterval int64
	sendInterval   int64
	addr           string
	hashKey        string
}

func New(logger logger.Logger, updateInterval int64, sendInterval int64, hashKey string, addr string) business.Agent {
	return &svc{
		logger: logger,

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},

		updateInterval: updateInterval,
		sendInterval:   sendInterval,
		addr:           addr,
		hashKey:        hashKey,
	}
}
