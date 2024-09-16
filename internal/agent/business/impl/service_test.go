package impl

import (
	"sync"
	"testing"

	semaphoreimpl "github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore/impl"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/mock"
)

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debugf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Infof(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Warnf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Errorf(format string, args ...interface{}) {
	m.Called(format, args)
}

func (m *MockLogger) Debugw(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Infow(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Warnw(msg string, keysAndValues ...interface{}) {
	m.Called(msg, keysAndValues)
}

func (m *MockLogger) Errorw(msg string, keysAndValues ...interface{}) {
	// m.Called(msg, keysAndValues)
}

func (m *MockLogger) Debug(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Info(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Warn(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Error(args ...interface{}) {
	m.Called(args)
}

func (m *MockLogger) Sync() error {
	return nil
}

func TestCollectMetrics(t *testing.T) {
	logger := MockLogger{}
	const (
		rateLimit = 10
	)
	s := &svc{
		logger: &logger,

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

func BenchmarkCollectMetrics(b *testing.B) {
	logger := MockLogger{}
	s := &svc{
		logger: &logger,

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
	logger := MockLogger{}
	s := &svc{
		logger: &logger,

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
