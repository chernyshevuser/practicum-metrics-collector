package impl

import (
	"context"
	"sync"
	"testing"
	"time"

	semaphoreimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore/impl"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"

	"github.com/shopspring/decimal"
)

func TestCollectMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)

	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}

	s.collectMetrics()

	cnt := len(s.metrics)
	if cnt <= 1 {
		t.Errorf("metrics len is invalid")
	}

	if s.pollCount.Cmp(decimal.NewFromInt(1)) != 0 {
		t.Errorf("expected pollCount to be 1, got %v", s.pollCount)
	}

	found := false
	for _, metric := range s.metrics {
		if metric.ID == "PollCount" {
			found = true
			if metric.Val.Cmp(decimal.NewFromInt(1)) != 0 {
				t.Errorf("expected PollCount metric value to be 1, got %v", metric.Val)
			}
		}
	}
	if !found {
		t.Errorf("PollCount metric not found")
	}

	s.collectMetrics()
	if s.pollCount.Cmp(decimal.NewFromInt(2)) != 0 {
		t.Errorf("expected pollCount to be 2, got %v", s.pollCount)
	}

	if cnt != len(s.metrics) {
		t.Errorf("metrics len is invalid")
	}
}

func TestCollectExtraMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)
	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}

	s.collectExtraMetrics()

	if len(s.metrics) != 0 {
		t.Errorf("metrics len is invalid")
	}

	if len(s.extraMetrics) == 0 {
		t.Errorf("extra metrics len is invalid")
	}

	for _, metric := range s.extraMetrics {
		if metric.ID == "PollCount" {
			t.Errorf("unexpected PollCount metric value: %v", metric.Val)
		}
	}
}

func TestClose(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)
	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}

	logger.EXPECT().Info("goodbye from agent-svc").Times(1)
	s.Close()
}

func TestRun_ctx(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)
	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},

		updateInterval: 100,
		sendInterval:   100,
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 100*time.Millisecond)
	defer cancel()

	logger.EXPECT().Infow("ctx done").Times(3)
	s.Run(ctx)

	time.Sleep(time.Second)
}

func TestRun_close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)
	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},

		updateInterval: 100,
		sendInterval:   100,
	}

	logger.EXPECT().Info("goodbye from agent-svc").Times(1)
	logger.EXPECT().Infow("close done").Times(3)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(300 * time.Millisecond)
		s.Close()
	}()

	s.Run(context.TODO())

	wg.Wait()
}
func TestRun_collectAllMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

	const (
		rateLimit = 10
	)
	s := &svc{
		logger: logger,

		semaphore: semaphoreimpl.New(rateLimit),

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},

		updateInterval: 1,
		sendInterval:   100,
	}

	logger.EXPECT().Infow("update extra metrics", "status", "start").Times(1)
	logger.EXPECT().Infow("update extra metrics", "status", "finished").Times(1)
	logger.EXPECT().Infow("update metrics", "status", "start").Times(1)
	logger.EXPECT().Infow("update metrics", "status", "finished").Times(1)
	logger.EXPECT().Info("goodbye from agent-svc").Times(1)
	logger.EXPECT().Infow("close done").Times(3)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(2 * time.Second)
		s.Close()
	}()

	s.Run(context.TODO())

	wg.Wait()
}

func BenchmarkCollectMetrics(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)
	s := &svc{
		logger: logger,

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.collectMetrics()
	}
}

func BenchmarkCollectExtraMetrics(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)
	s := &svc{
		logger: logger,

		pollCount: decimal.Decimal{},
		mu:        &sync.Mutex{},

		closeCh: make(chan struct{}, 1),
		wg:      &sync.WaitGroup{},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		s.collectExtraMetrics()
	}
}
