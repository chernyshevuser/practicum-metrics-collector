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

// TODO check method GET/POST
func SetupRouter(api handler.API, router *mux.Router, logger logger.Logger) {
	router.HandleFunc(UpdateMetricPath, middleware.Accept(api.UpdateMetric, logger)).Methods(http.MethodPost)
	router.HandleFunc(UpdateMetricJSONPath, middleware.Accept(api.UpdateMetricJSON, logger)).Methods(http.MethodPost)
	router.HandleFunc(UpdateMetricsJSONPath, middleware.Accept(api.UpdateMetricsJSON, logger)).Methods(http.MethodPost)
	router.HandleFunc(GetMetricValuePath, middleware.Accept(api.GetMetricValue, logger)).Methods(http.MethodGet)
	router.HandleFunc(GetMetricValueJSONPath, middleware.Accept(api.GetMetricValueJSON, logger)).Methods(http.MethodGet)
	router.HandleFunc(GetAllMetricsPath, middleware.Accept(api.GetAllMetrics, logger)).Methods(http.MethodGet)
	router.HandleFunc(PingDB, middleware.Accept(api.PingDB, logger)).Methods(http.MethodGet)
}
