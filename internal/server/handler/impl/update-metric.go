package impl

import (
	"net/http"
	"strings"
)

func (a *api) UpdateMetric(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	status := http.StatusOK

	ctx := r.Context()

	vals := strings.Split(r.URL.String(), "/")
	if len(vals) != 5 {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	metricTypeStr := vals[2]
	metricNameStr := vals[3]
	metricValueStr := vals[4]

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
