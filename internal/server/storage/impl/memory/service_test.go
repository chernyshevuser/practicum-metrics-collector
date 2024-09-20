package memorystorage_test

import (
	"context"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	memorystorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/impl/memory"
	logger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/assert"
)

func TestSetGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info("goodbye from db-svc").Times(1)

	fname, flag := "", false
	svc, err := memorystorage.New(context.TODO(), mockLogger, fname, flag)
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
		err = svc.Set(context.TODO(), m)
		assert.NoError(t, err)
		svc.Unlock()
	}

	for _, m := range metrics {
		svc.Lock()
		retrievedMetric, err := svc.Get(context.TODO(), storage.BuildKey(m.ID, m.Type))
		assert.NoError(t, err)
		assert.Equal(t, m, *retrievedMetric)
		svc.Unlock()
	}
}

func TestGetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info("goodbye from db-svc").Times(1)

	fname, flag := "", false
	svc, err := memorystorage.New(context.TODO(), mockLogger, fname, flag)
	if err != nil {
		t.Errorf("can't create storage svc: %v", err)
	}
	defer svc.Close()

	metric1 := storage.Metric{ID: "id1", Type: "counter", Val: 10}
	metric2 := storage.Metric{ID: "id2", Type: "gauge", Val: 5.5}
	svc.Lock()
	svc.Set(context.TODO(), metric1)
	svc.Set(context.TODO(), metric2)
	svc.Unlock()

	svc.Lock()
	allMetrics, err := svc.GetAll(context.TODO())
	svc.Unlock()
	assert.NoError(t, err)
	assert.Len(t, *allMetrics, 2)
	assert.Contains(t, *allMetrics, metric1)
	assert.Contains(t, *allMetrics, metric2)
}

func BenchmarkSet(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info("goodbye from db-svc").Times(1)

	fname, flag := "", false
	svc, err := memorystorage.New(context.TODO(), mockLogger, fname, flag)
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
		err := svc.Set(context.TODO(), metric)
		if err != nil {
			b.Fatalf("set failed: %v", err)
		}
		svc.Unlock()
	}
}

func BenchmarkGet(b *testing.B) {
	ctrl := gomock.NewController(b)
	defer ctrl.Finish()

	mockLogger := logger.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info("goodbye from db-svc").Times(1)

	fname, flag := "", false
	svc, err := memorystorage.New(context.TODO(), mockLogger, fname, flag)
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
	svc.Set(context.TODO(), metric)
	svc.Unlock()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		svc.Lock()
		_, err := svc.Get(context.TODO(), storage.BuildKey(metric.ID, metric.Type))
		if err != nil {
			b.Fatalf("get failed: %v", err)
		}
		svc.Unlock()
	}
}
