package memorystorage_test

import (
	"context"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	memorystorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/memory"
	"github.com/test-go/testify/assert"
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
	m.Called(msg, keysAndValues)
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
	return m.Called().Error(0)
}

func TestSetGet(t *testing.T) {
	ctx := context.TODO()
	mockLogger := new(MockLogger)
	mockLogger.On("Info", []interface{}{"goodbye from db-svc"}).Return()
	fname, flag := "", false

	svc, err := memorystorage.New(ctx, mockLogger, fname, flag)
	if err != nil {
		t.Errorf("can't create storage svc: %v", err)
	}
	defer svc.Close()

	metrics := []storage.Metric{
		{
			ID:    "id",
			Type:  "counter",
			Delta: 100,
		},
		{
			ID:   "1",
			Type: "gauge",
			Val:  42.0,
		},
	}

	for _, m := range metrics {
		svc.Lock()
		err = svc.Set(ctx, m)
		assert.NoError(t, err)
		svc.Unlock()
	}

	for _, m := range metrics {
		svc.Lock()
		retrievedMetric, err := svc.Get(ctx, storage.BuildKey(m.ID, m.Type))
		assert.NoError(t, err)
		assert.Equal(t, m, *retrievedMetric)
		svc.Unlock()
	}
}

func TestGetAll(t *testing.T) {
	ctx := context.TODO()
	mockLogger := new(MockLogger)
	mockLogger.On("Info", []interface{}{"goodbye from db-svc"}).Return()
	fname, flag := "", false

	svc, err := memorystorage.New(ctx, mockLogger, fname, flag)
	if err != nil {
		t.Errorf("can't create storage svc: %v", err)
	}
	defer svc.Close()

	metric1 := storage.Metric{ID: "id1", Type: "counter", Val: 10}
	metric2 := storage.Metric{ID: "id2", Type: "gauge", Val: 5.5}
	svc.Lock()
	svc.Set(ctx, metric1)
	svc.Set(ctx, metric2)
	svc.Unlock()

	svc.Lock()
	allMetrics, err := svc.GetAll(ctx)
	svc.Unlock()
	assert.NoError(t, err)
	assert.Len(t, *allMetrics, 2)
	assert.Contains(t, *allMetrics, metric1)
	assert.Contains(t, *allMetrics, metric2)
}

func BenchmarkSet(b *testing.B) {
	ctx := context.TODO()
	mockLogger := new(MockLogger)
	mockLogger.On("Info", []interface{}{"goodbye from db-svc"}).Return()
	fname, flag := "", false

	svc, err := memorystorage.New(ctx, mockLogger, fname, flag)
	if err != nil {
		b.Errorf("can't create storage svc: %v", err)
	}
	defer svc.Close()

	metric := storage.Metric{
		ID:   "benchmark-id",
		Type: "gauge",
		Val:  42.0,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.Lock()
		err := svc.Set(ctx, metric)
		if err != nil {
			b.Fatalf("set failed: %v", err)
		}
		svc.Unlock()
	}
}

func BenchmarkGet(b *testing.B) {
	ctx := context.TODO()
	mockLogger := new(MockLogger)
	mockLogger.On("Info", []interface{}{"goodbye from db-svc"}).Return()
	fname, flag := "", false

	svc, err := memorystorage.New(ctx, mockLogger, fname, flag)
	if err != nil {
		b.Errorf("can't create storage svc: %v", err)
	}
	defer svc.Close()

	metric := storage.Metric{
		ID:   "benchmark-id",
		Type: "gauge",
		Val:  42.0,
	}
	svc.Lock()
	svc.Set(ctx, metric)
	svc.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.Lock()
		_, err := svc.Get(ctx, storage.BuildKey(metric.ID, metric.Type))
		if err != nil {
			b.Fatalf("get failed: %v", err)
		}
		svc.Unlock()
	}
}
