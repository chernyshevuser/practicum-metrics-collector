package handler

import "net/http"

type API interface {
	UpdateMetric(w http.ResponseWriter, r *http.Request) error
	UpdateMetricJSON(w http.ResponseWriter, r *http.Request) error
	UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) error
	GetMetricValue(w http.ResponseWriter, r *http.Request) error
	GetMetricValueJSON(w http.ResponseWriter, r *http.Request) error
	GetAllMetrics(w http.ResponseWriter, r *http.Request) error
	PingDB(w http.ResponseWriter, r *http.Request) error
}
