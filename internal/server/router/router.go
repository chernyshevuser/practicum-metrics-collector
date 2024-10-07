package router

import (
	"net/http"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/handler"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/middleware"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/logger"
	"github.com/gorilla/mux"
)

const (
	UpdateMetricPath       string = "/update/{type}/{name}/{value}"
	UpdateMetricJSONPath   string = "/update/"
	UpdateMetricsJSONPath  string = "/updates/"
	GetMetricValuePath     string = "/value/{type}/{name}"
	GetMetricValueJSONPath string = "/value/"
	GetAllMetricsPath      string = "/"
	PingDB                 string = "/ping"
)

func SetupRouter(api handler.API, router *mux.Router, logger logger.Logger, cryptoKey string) {
	router.HandleFunc(UpdateMetricPath, middleware.Accept(api.UpdateMetric, logger, cryptoKey)).Methods(http.MethodPost)
	router.HandleFunc(UpdateMetricJSONPath, middleware.Accept(api.UpdateMetricJSON, logger, cryptoKey)).Methods(http.MethodPost)
	router.HandleFunc(UpdateMetricsJSONPath, middleware.Accept(api.UpdateMetricsJSON, logger, cryptoKey)).Methods(http.MethodPost)
	router.HandleFunc(GetMetricValuePath, middleware.Accept(api.GetMetricValue, logger, cryptoKey)).Methods(http.MethodGet)
	router.HandleFunc(GetMetricValueJSONPath, middleware.Accept(api.GetMetricValueJSON, logger, cryptoKey)).Methods(http.MethodPost)
	router.HandleFunc(GetAllMetricsPath, middleware.Accept(api.GetAllMetrics, logger, cryptoKey)).Methods(http.MethodGet)
	router.HandleFunc(PingDB, middleware.Accept(api.PingDB, logger, cryptoKey)).Methods(http.MethodGet)
}
