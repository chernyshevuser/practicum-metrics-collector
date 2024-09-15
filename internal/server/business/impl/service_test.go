package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
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

type MockBusinessSvc struct {
	mock.Mock
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) Set(ctx context.Context, metric storage.Metric) error {
	args := m.Called(ctx, metric)
	return args.Error(0)
}

func (m *MockStorage) Get(ctx context.Context, key string) (*storage.Metric, error) {
	args := m.Called(ctx, key)
	if metric, ok := args.Get(0).(*storage.Metric); ok {
		return metric, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) GetAll(ctx context.Context) (*[]storage.Metric, error) {
	args := m.Called(ctx)
	if metrics, ok := args.Get(0).(*[]storage.Metric); ok {
		return metrics, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockStorage) Lock() {
	m.Called()
}

func (m *MockStorage) Unlock() {
	m.Called()
}

func (m *MockStorage) Actualize(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) Dump(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) Ping(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockStorage) Close() error {
	args := m.Called()
	return args.Error(0)
}

func TestGetMetricValue(t *testing.T) {
	ctx := context.TODO()
	st := &MockStorage{}
	logger := &MockLogger{}

	svc := impl.New(st, logger)

	tests := []struct {
		name               string
		metricType         string
		metricName         string
		mockSetup          func()
		expectedValue      decimal.Decimal
		expectedMetricType business.MetricType
		expectError        bool
	}{
		{
			name:       "counter metric found",
			metricType: "counter",
			metricName: "sample_counter",
			mockSetup: func() {
				st.On("Lock").Return()
				st.On("Unlock").Return()
				st.On("Get", mock.Anything, storage.BuildKey("sample_counter", "counter")).Return(&storage.Metric{Delta: 123}, nil)
			},
			expectedValue:      decimal.NewFromInt(123),
			expectedMetricType: business.Counter,
			expectError:        false,
		},
		{
			name:       "gauge metric found",
			metricType: "gauge",
			metricName: "sample_gauge",
			mockSetup: func() {
				st.On("Lock").Return()
				st.On("Unlock").Return()
				st.On("Get", mock.Anything, storage.BuildKey("sample_gauge", "gauge")).Return(&storage.Metric{Val: 123.45}, nil)
			},
			expectedValue:      decimal.NewFromFloat(123.45),
			expectedMetricType: business.Gauge,
			expectError:        false,
		},
		{
			name:       "unknown metric type",
			metricType: "unknown",
			metricName: "invalid_metric",
			mockSetup: func() {
				logger.On("Infow", "parse metric", "given type", "unknown", "parsed type", business.Unknown, "given name", "invalid_metric").Return()
			},
			expectedValue:      decimal.Decimal{},
			expectedMetricType: "",
			expectError:        true,
		},
		{
			name:       "error fetching from storage",
			metricType: "counter",
			metricName: "error_metric",
			mockSetup: func() {
				st.On("Lock").Return()
				st.On("Unlock").Return()

				st.On("Get", mock.Anything, storage.BuildKey("error_metric", "counter")).Return(nil, fmt.Errorf("db error"))

				logger.On("Infow", "parse metric", "given type", "counter", "parsed type", business.Counter, "given name", "error_metric").Return()
				logger.On("Errorw", "storage problem", "msg", "can't get counter metric val", "reason", mock.Anything).Return()
			},
			expectedValue:      decimal.Decimal{},
			expectedMetricType: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			val, metricType, err := svc.GetMetricValue(ctx, tt.metricType, tt.metricName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, *val)
				assert.Equal(t, tt.expectedMetricType, metricType)
			}

			st.AssertExpectations(t)
		})
	}
}
