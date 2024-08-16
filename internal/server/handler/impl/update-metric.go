package impl

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
)

func (a *api) UpdateMetric(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	status := http.StatusOK

	ctx := r.Context()

	// can't be used cause of mux can't pass Vars inside the test
	// vars := mux.Vars(r)
	// metricTypeStr := vars["type"]
	// metricNameStr := vars["name"]
	// metricValueStr := vars["value"]

	values := strings.Split(r.URL.String(), "/")
	metricTypeStr := values[2]
	metricNameStr := values[3]
	metricValueStr := values[4]

	if len(metricTypeStr) == 0 || len(metricNameStr) == 0 || len(metricValueStr) == 0 {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	err := a.mc.UpdateMetrics(ctx, []business.RawMetric{
		{
			ID:    metricNameStr,
			Type:  metricTypeStr,
			Value: metricValueStr,
		},
	})
	if err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	w.WriteHeader(status)
	return nil
}

type updateMetricReq struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type updateMetricResp struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (a *api) UpdateMetricJSON(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK

	ctx := r.Context()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	var req updateMetricReq
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
	}

	metricNameStr := req.ID
	metricTypeStr := req.MType

	var metricValueStr string

	if req.Delta == nil && req.Value == nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	} else if req.Delta != nil {
		metricValueStr = fmt.Sprintf("%v", *req.Delta)
	} else if req.Value != nil {
		metricValueStr = fmt.Sprintf("%v", *req.Value)
	}

	err := a.mc.UpdateMetrics(ctx, []business.RawMetric{
		{
			ID:    metricNameStr,
			Type:  metricTypeStr,
			Value: metricValueStr,
		},
	})
	if err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	w.WriteHeader(status)

	resp := updateMetricResp(req)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err = w.Write(respBytes); err != nil {
		return err
	}

	return nil
}

type updateMetricsReq struct {
	Metrics []struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}
}

type updateMetricsResp struct {
	Metrics []struct {
		ID    string   `json:"id"`
		MType string   `json:"type"`
		Delta *int64   `json:"delta,omitempty"`
		Value *float64 `json:"value,omitempty"`
	}
}

func (a *api) UpdateMetricsJSON(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")
	status := http.StatusOK

	ctx := r.Context()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r.Body); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	var req updateMetricsReq
	if err := json.Unmarshal(buf.Bytes(), &req); err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
	}

	var rawMetrics []business.RawMetric

	for _, m := range req.Metrics {
		metricNameStr := m.ID
		metricTypeStr := m.MType

		var metricValueStr string

		if m.Delta == nil && m.Value == nil {
			status = http.StatusBadRequest
			w.WriteHeader(status)
			return nil
		} else if m.Delta != nil {
			metricValueStr = fmt.Sprintf("%v", *m.Delta)
		} else if m.Value != nil {
			metricValueStr = fmt.Sprintf("%v", *m.Value)
		}

		rawMetrics = append(rawMetrics, business.RawMetric{
			ID:    metricNameStr,
			Type:  metricTypeStr,
			Value: metricValueStr,
		})
	}

	err := a.mc.UpdateMetrics(ctx, rawMetrics)
	if err != nil {
		status = http.StatusBadRequest
		w.WriteHeader(status)
		return nil
	}

	w.WriteHeader(status)

	resp := updateMetricsResp(req)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	if _, err = w.Write(respBytes); err != nil {
		return err
	}

	return nil
}
