package impl

import (
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore"
	semaphoreimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore/impl"
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

	semaphore semaphore.Semaphore

	metrics      []Metric
	extraMetrics []Metric
	pollCount    decimal.Decimal
	mu           *sync.Mutex

	closeCh chan struct{}
	wg      *sync.WaitGroup

	updateInterval int64
	sendInterval   int64
	addr           string
	hashKey        string

	cryptoKey string
}

func New(logger logger.Logger, updateInterval int64, sendInterval int64, hashKey string, addr string, rateLimit int64, cryptoKey string) business.Agent {
	return &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},

		updateInterval: updateInterval,
		sendInterval:   sendInterval,
		addr:           addr,
		hashKey:        hashKey,
		cryptoKey:      cryptoKey,
	}
}
