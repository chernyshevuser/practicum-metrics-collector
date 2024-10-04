package middleware_test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/middleware"
	mocklogger "github.com/chernyshevuser/practicum-metrics-collector/tools/logger/mock"
	"github.com/golang/mock/gomock"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func TestLogMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	logger.EXPECT().Infow("entering", "method", "GET", "url", "/test").Times(1)
	logger.EXPECT().Infow("leaving", "method", "GET", "url", "/test").Times(1)

	handler := middleware.LogMiddleware(testHandler, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestPanicMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	logger.EXPECT().Errorw("panic happened", "reason", "test panic", "stacktrace", gomock.Any()).Times(1)

	handler := middleware.PanicMiddleware(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}

func TestCompressMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	handler := middleware.CompressMiddleware(testHandler, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.Header.Get("Content-Encoding") != "gzip" {
		t.Errorf("expected gzip encoding")
	}

	gzReader, err := gzip.NewReader(resp.Body)
	if err != nil {
		t.Fatalf("failed to create gzip reader: %v", err)
	}
	body, _ := io.ReadAll(gzReader)

	if string(body) != "OK" {
		t.Errorf("expected body 'OK', got '%s'", body)
	}
}

func TestDecompressMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	handler := middleware.DecompressMiddleware(testHandler, logger)

	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	gzWriter.Write([]byte("gzip request"))
	gzWriter.Close()

	req := httptest.NewRequest(http.MethodPost, "/test", &buf)
	req.Header.Set("Content-Encoding", "gzip")
	w := httptest.NewRecorder()

	handler(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if string(body) != "OK" {
		t.Errorf("expected body 'OK', got '%s'", body)
	}
}

func TestErrorMiddleware(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := mocklogger.NewMockLogger(ctrl)
	logger.EXPECT().Errorw(
		"error happened",
		"url", "/test",
		"reason", gomock.Any(),
	).Times(1)

	handler := middleware.ErrorMiddleware(func(w http.ResponseWriter, r *http.Request) error {
		return fmt.Errorf("test error")
	}, logger)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	handler(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", w.Code)
	}
}
