package impl

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	mockbusiness "github.com/chernyshevuser/practicum-metrics-collector/internal/server/business/mock"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/router"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
)

// Example_api_GetMetricValueJSON demonstrates how to use GetMetricValueJSON handler.
func Example_api_GetMetricValueJSON() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)

	val := decimal.NewFromInt(123)
	businessSvc.EXPECT().GetMetricValue(gomock.Any(), "counter", "some_id").Return(&val, business.Counter, nil)

	svc := New(businessSvc, logger)

	reqBody := `{
		"id": "some_id",
		"type": "counter"
	}`
	req := httptest.NewRequest(http.MethodPost, router.GetMetricValuePath, bytes.NewBuffer([]byte(reqBody)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	svc.GetMetricValueJSON(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// {"id":"some_id","type":"counter","delta":123}
}

// Example_api_GetAllMetrics demonstrates how to use GetAllMetrics handler.
func Example_api_GetAllMetrics() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)

	businessSvc.EXPECT().GetAllMetrics(gomock.Any()).Return([]business.CounterMetric{}, []business.GaugeMetric{}, nil)

	svc := New(businessSvc, logger)
	_ = svc

	req := httptest.NewRequest(http.MethodGet, router.GetAllMetricsPath, nil)
	w := httptest.NewRecorder()

	svc.GetAllMetrics(w, req)

	res := w.Result()
	defer res.Body.Close()
	fmt.Println(res.StatusCode)

	body, _ := io.ReadAll(res.Body)
	fmt.Println(string(body))

	// Output:
	// 200
	// [<html><body><h1></h1><ul> </ul></body></html>]
}

// Example_api_GetMetricValue demonstrates how to use GetMetricValue handler.
func Example_api_GetMetricValue() {
	ctrl := gomock.NewController(nil)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	businessSvc := mockbusiness.NewMockMetricsCollector(ctrl)
	val := decimal.NewFromInt(123)
	businessSvc.EXPECT().GetMetricValue(gomock.Any(), "counter", "some_name").Return(&val, business.Counter, nil)

	svc := New(businessSvc, logger)

	req := httptest.NewRequest(http.MethodGet, router.GetMetricValuePath, nil)
	req = mux.SetURLVars(req, map[string]string{
		"type": "counter",
		"name": "some_name",
	})
	w := httptest.NewRecorder()

	svc.GetMetricValue(w, req)

	res := w.Result()
	defer res.Body.Close()

	fmt.Println(res.StatusCode)

	// Output: 200
}
