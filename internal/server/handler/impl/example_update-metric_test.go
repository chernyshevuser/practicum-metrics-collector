package impl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	mockbusiness "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/mock"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler/impl"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/router"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
)

// ExampleUpdateMetricJSON demonstrates how to use UpdateMetricJSON handler.
func ExampleUpdateMetricJSON() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)
	businessSvc.EXPECT().UpdateMetrics(
		gomock.Any(),
		[]business.RawMetric{
			{
				ID:    "some_name",
				Type:  "counter",
				Value: fmt.Sprintf("%v", decimal.NewFromInt(123)),
			},
		}).Return(
		[]business.CounterMetric{
			{
				ID:    "some_name",
				Delta: decimal.NewFromInt(123),
			},
		},
		[]business.GaugeMetric{},
		nil,
	)

	svc := impl.New(businessSvc, logger)

	val := int64(123)
	reqBody, _ := json.Marshal(struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}{
		ID:    "some_name",
		MType: "counter",
		Delta: &val,
	})

	req := httptest.NewRequest(http.MethodPost, router.UpdateMetricJSONPath, bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.UpdateMetricJSON(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}

// ExampleUpdateMetricsJSON demonstrates how to use UpdateMetricsJSON handler.
func ExampleUpdateMetricsJSON() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)
	businessSvc.EXPECT().UpdateMetrics(
		gomock.Any(),
		[]business.RawMetric{
			{
				ID:    "some_name",
				Type:  "counter",
				Value: fmt.Sprintf("%v", decimal.NewFromInt(123)),
			},
		}).Return(
		[]business.CounterMetric{
			{
				ID:    "some_name",
				Delta: decimal.NewFromInt(123),
			},
		},
		[]business.GaugeMetric{},
		nil,
	)

	svc := impl.New(businessSvc, logger)

	val := int64(123)
	reqBody, _ := json.Marshal([]struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}{
		{
			ID:    "some_name",
			MType: "counter",
			Delta: &val,
		},
	})

	req := httptest.NewRequest(http.MethodPost, router.UpdateMetricsJSONPath, bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.UpdateMetricsJSON(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output:
	// 200
}
