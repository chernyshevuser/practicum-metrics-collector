package impl

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *api) GetMetricValue(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	status := http.StatusOK

	ctx := r.Context()

	vars := mux.Vars(r)
	metricTypeStr := vars["type"]
	metricNameStr := vars["name"]

	if len(metricTypeStr) == 0 || len(metricNameStr) == 0 {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	val, err := a.mc.GetMetricValue(ctx, metricTypeStr, metricNameStr)
	if err != nil {
		return fmt.Errorf("can't get metric val, reason: %v", err)
	}

	if val == nil {
		status = http.StatusNotFound
		w.WriteHeader(status)
		return nil
	}

	_, err = fmt.Fprint(w, val.String())
	if err != nil {
		return fmt.Errorf("can't write metric val to responseWriter, reason: %v", err)
	}

	w.WriteHeader(status)
	return nil
}

func (a *api) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) error {
	panic("")
}

func (a *api) GetAllMetrics(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/html")
	status := http.StatusOK

	ctx := r.Context()

	counterMetrics, gaugeMetrics, err := a.mc.GetAllMetrics(ctx)
	if err != nil {
		return err
	}

	htmlRes := []string{"<html><body><h1></h1><ul>"}

	for _, m := range counterMetrics {
		htmlRes = append(htmlRes, fmt.Sprintf("<li>%s: %s</li>", m.ID, m.Delta.String()))
	}

	for _, m := range gaugeMetrics {
		htmlRes = append(htmlRes, fmt.Sprintf("<li>%s: %s</li>", m.ID, m.Value.String()))
	}

	htmlRes = append(htmlRes, "</ul></body></html>")

	//should be written before Fprint(w, ...)
	w.WriteHeader(status)

	_, err = fmt.Fprint(w, htmlRes)
	if err != nil {
		return err
	}

	return nil
}
