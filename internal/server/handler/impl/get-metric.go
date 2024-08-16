package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
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

	val, _, err := a.mc.GetMetricValue(ctx, metricTypeStr, metricNameStr)
	if err != nil {
		return fmt.Errorf("can't get metric val, reason: %v", err)
	}

	if val == nil {
		status = http.StatusNotFound
		w.WriteHeader(status)
		return nil
	}

	w.WriteHeader(status)

	_, err = fmt.Fprint(w, val.String())
	if err != nil {
		return fmt.Errorf("can't write metric val to responseWriter, reason: %v", err)
	}

	return nil
}

type getMetricReq struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type getMetricResp struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (a *api) GetMetricValueJSON(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	status := http.StatusOK

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	var req getMetricReq
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	val, mType, err := a.mc.GetMetricValue(ctx, req.MType, req.ID)
	if err != nil {
		return fmt.Errorf("can't get metric val, reason: %v", err)
	}

	if val == nil {
		status = http.StatusNotFound
		w.WriteHeader(status)
		return nil
	}

	resp := getMetricResp{
		ID:    req.ID,
		MType: req.MType,
	}

	if mType == business.Counter {
		valInt64 := val.IntPart()
		resp.Delta = &valInt64
	} else if mType == business.Gauge {
		valFloat64 := val.InexactFloat64()
		resp.Value = &valFloat64
	} else {
		a.logger.Errorw(
			"unknown metric type",
			"type", mType,
		)
		return fmt.Errorf("get unknown metric type from business")
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err = w.Write(respBytes); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
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
