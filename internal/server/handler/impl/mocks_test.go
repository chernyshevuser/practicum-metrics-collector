package impl_test

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
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
	return nil
}

type MockBusinessSvc struct {
	mock.Mock
}

func (m *MockBusinessSvc) UpdateMetrics(ctx context.Context, metrics []business.RawMetric) (updatedCounterMetric []business.CounterMetric, updatedGaugeMetrics []business.GaugeMetric, err error) {
	args := m.Called(ctx, metrics)
	if arg := args.Get(0); arg != nil {
		updatedCounterMetric = arg.([]business.CounterMetric)
	}
	if arg := args.Get(1); arg != nil {
		updatedGaugeMetrics = arg.([]business.GaugeMetric)
	}
	err = args.Error(2)
	return
}

func (m *MockBusinessSvc) GetMetricValue(ctx context.Context, metricType, metricName string) (val *decimal.Decimal, mType business.MetricType, err error) {
	args := m.Called(ctx, metricType, metricName)
	if arg := args.Get(0); arg != nil {
		val = arg.(*decimal.Decimal)
	}
	if arg := args.Get(1); arg != nil {
		mType = arg.(business.MetricType)
	}
	err = args.Error(2)
	return
}

func (m *MockBusinessSvc) GetAllMetrics(ctx context.Context) (counterMetrics []business.CounterMetric, gaugeMetrics []business.GaugeMetric, err error) {
	args := m.Called(ctx)
	if arg := args.Get(0); arg != nil {
		counterMetrics = arg.([]business.CounterMetric)
	}
	if arg := args.Get(1); arg != nil {
		gaugeMetrics = arg.([]business.GaugeMetric)
	}
	err = args.Error(2)
	return
}

func (m *MockBusinessSvc) PingDB(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockBusinessSvc) Close() {}
