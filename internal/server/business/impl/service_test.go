package impl_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	mockstorage "github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage/mock"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

type MockBusinessSvc struct {
	mock.Mock
}

func TestGetMetricValue(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := mockstorage.NewMockStorage(ctrl)
	defer st.Close()
	st.EXPECT().Close().Times(1)

	logger := mocklogger.NewMockLogger(ctrl)
	defer logger.Sync()
	logger.EXPECT().Sync().Times(1)

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
				st.EXPECT().Lock().Times(1)
				st.EXPECT().Unlock().Times(1)
				st.EXPECT().Get(gomock.Any(), storage.BuildKey("sample_counter", "counter")).Return(&storage.Metric{Delta: 123}, nil)
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
				st.EXPECT().Lock().Times(1)
				st.EXPECT().Unlock().Times(1)
				st.EXPECT().Get(gomock.Any(), storage.BuildKey("sample_gauge", "gauge")).Times(1).Return(&storage.Metric{Val: 123.45}, nil)
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
				st.EXPECT().Lock().Times(1)
				st.EXPECT().Unlock().Times(1)
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
				st.EXPECT().Lock().Times(1)
				st.EXPECT().Unlock().Times(1)
				st.EXPECT().Get(gomock.Any(), storage.BuildKey("error_metric", "counter")).Return(nil, fmt.Errorf("db error")).Times(1)
				logger.EXPECT().Errorw("storage problem", "msg", "can't get counter metric val", "reason", gomock.Any()).Times(1)
			},
			expectedValue:      decimal.Decimal{},
			expectedMetricType: "",
			expectError:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			val, metricType, err := svc.GetMetricValue(context.TODO(), tt.metricType, tt.metricName)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedValue, *val)
				assert.Equal(t, tt.expectedMetricType, metricType)
			}
		})
	}
}
