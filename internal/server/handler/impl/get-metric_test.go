package impl_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler/impl"
	"github.com/gorilla/mux"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/assert"
	"github.com/test-go/testify/mock"
)

func TestGetMetricValue(t *testing.T) {
	businessSvc := &MockBusinessSvc{}
	defer businessSvc.Close()
	logger := &MockLogger{}
	defer logger.Sync()

	svc := impl.New(businessSvc, logger)

	type want struct {
		code        int
		contentType string
		body        string
	}

	tests := []struct {
		name     string
		endp     string
		mtype    string
		mname    string
		mockResp func()
		want     want
	}{
		{
			name:  "valid request with metric value",
			endp:  "/metric/counter/sample_text",
			mtype: "counter",
			mname: "sample_text",
			mockResp: func() {
				val := decimal.NewFromInt(123)
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "sample_text").
					Return(&val, business.Counter, nil)
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				body:        "123",
			},
		},
		{
			name:  "metric not found",
			endp:  "/metric/counter/sample_text_not_found",
			mtype: "counter",
			mname: "sample_text_not_found",
			mockResp: func() {
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "sample_text_not_found").
					Return(nil, business.Counter, nil)
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "text/plain; charset=utf-8",
				body:        "",
			},
		},
		{
			name:  "error from service",
			endp:  "/metric/counter/error_from_svc",
			mtype: "counter",
			mname: "error_from_svc",
			mockResp: func() {
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "error_from_svc").
					Return(nil, business.Counter, fmt.Errorf("mock error"))
			},
			want: want{
				code:        http.StatusOK,
				contentType: "text/plain; charset=utf-8",
				body:        "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockResp()
			req := httptest.NewRequest(http.MethodGet, test.endp, nil)
			vars := map[string]string{
				"type": test.mtype,
				"name": test.mname,
			}

			req = mux.SetURLVars(req, vars)
			w := httptest.NewRecorder()
			err := svc.GetMetricValue(w, req)

			res := w.Result()
			defer res.Body.Close()

			if err != nil && !strings.Contains(err.Error(), "mock") {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.code, res.StatusCode)
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.want.body, string(body))
		})
	}
}

func TestGetMetricValueJSON(t *testing.T) {
	businessSvc := &MockBusinessSvc{}
	defer businessSvc.Close()
	logger := &MockLogger{}
	defer logger.Sync()

	svc := impl.New(businessSvc, logger)

	type want struct {
		code        int
		contentType string
		body        string
	}

	type getMetricReq struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}

	tests := []struct {
		name     string
		request  interface{}
		mockResp func()
		want     want
	}{
		{
			name: "valid counter metric",
			request: getMetricReq{
				ID:    "sample_text",
				MType: "counter",
			},
			mockResp: func() {
				val := decimal.NewFromInt(123)
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "sample_text").
					Return(&val, business.Counter, nil)
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"sample_text","type":"counter","delta":123}`,
			},
		},
		{
			name: "valid gauge metric",
			request: getMetricReq{
				ID:    "sample_text",
				MType: "gauge",
			},
			mockResp: func() {
				val := decimal.NewFromFloat(123.45)
				businessSvc.On("GetMetricValue", mock.Anything, "gauge", "sample_text").
					Return(&val, business.Gauge, nil)
			},
			want: want{
				code:        http.StatusOK,
				contentType: "application/json",
				body:        `{"id":"sample_text","type":"gauge","value":123.45}`,
			},
		},
		{
			name: "metric not found",
			request: getMetricReq{
				ID:    "sample_text_not_found",
				MType: "counter",
			},
			mockResp: func() {
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "sample_text_not_found").
					Return(nil, business.Counter, nil)
			},
			want: want{
				code:        http.StatusNotFound,
				contentType: "application/json",
				body:        "",
			},
		},
		{
			name: "error from service",
			request: getMetricReq{
				ID:    "sample_text_err_from_svc",
				MType: "counter",
			},
			mockResp: func() {
				businessSvc.On("GetMetricValue", mock.Anything, "counter", "sample_text_err_from_svc").
					Return(nil, business.Counter, fmt.Errorf("mock error"))
			},
			want: want{
				code:        http.StatusOK, // changes in middleware
				contentType: "application/json",
				body:        "",
			},
		},
		{
			name: "unknown metric type",
			request: getMetricReq{
				ID:    "sample_text",
				MType: "unknown",
			},
			mockResp: func() {
				val := decimal.NewFromInt(123)
				businessSvc.On("GetMetricValue", mock.Anything, "unknown", "sample_text").
					Return(&val, business.Unknown, fmt.Errorf("mock error"))
			},
			want: want{
				code:        http.StatusOK, // changes in middleware
				contentType: "application/json",
				body:        "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.request)
			if err != nil {
				t.Fatalf("could not encode request: %v", err)
			}

			test.mockResp()
			request := httptest.NewRequest(http.MethodPost, "/get-metric-value", &buf)
			w := httptest.NewRecorder()

			err = svc.GetMetricValueJSON(w, request)

			res := w.Result()
			defer res.Body.Close()

			if err != nil && !strings.Contains(err.Error(), "mock") {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Equal(t, test.want.contentType, res.Header.Get("Content-Type"))
			assert.Equal(t, test.want.code, res.StatusCode)
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, test.want.body, string(body))
		})
	}
}
