package impl

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (a *api) UpdateMetric(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	status := http.StatusOK

	ctx := r.Context()

	vars := mux.Vars(r)
	metricTypeStr := vars["type"]
	metricNameStr := vars["name"]
	metricValueStr := vars["value"]

	if len(metricTypeStr) == 0 || len(metricNameStr) == 0 || len(metricValueStr) == 0 {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	err := a.mc.UpdateMetric(ctx, metricTypeStr, metricNameStr, metricValueStr)
	if err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	w.WriteHeader(status)
	return nil
}

func (a *api) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) error {
	panic("")
}

func (a *api) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) error {
	panic("")
}
