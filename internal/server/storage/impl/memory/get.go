package memorystorage

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
)

func (s *svc) GetGauge(ctx context.Context, id string) (*storage.GaugeMetric, error) {
	val, ok := s.gaugeStorage.Get(id)
	if !ok {
		return nil, nil
	}

	if val == nil {
		return nil, nil
	}

	m, ok := val.(storage.GaugeMetric)
	if !ok {
		return nil, fmt.Errorf("value has wrong type")
	}

	return &m, nil
}

func (s *svc) GetCounter(ctx context.Context, id string) (*storage.CounterMetric, error) {
	val, ok := s.counterStorage.Get(id)
	if !ok {
		return nil, nil
	}

	if val == nil {
		return nil, nil
	}

	m, ok := val.(storage.CounterMetric)
	if !ok {
		return nil, fmt.Errorf("value has wrong type")
	}

	return &m, nil
}

func (s *svc) GetAll(ctx context.Context) (*[]storage.GaugeMetric, *[]storage.CounterMetric, error) {
	counterData := s.counterStorage.GetAll()
	gaugeData := s.gaugeStorage.GetAll()

	counterMetrics := make([]storage.CounterMetric, 0, len(counterData))
	gaugeMetrics := make([]storage.GaugeMetric, 0, len(gaugeData))

	for _, tmp := range gaugeData {
		gaugeMetrics = append(gaugeMetrics, storage.GaugeMetric{
			ID:    tmp.(storage.GaugeMetric).ID,
			Value: tmp.(storage.GaugeMetric).Value,
		})
	}

	for _, tmp := range counterData {
		counterMetrics = append(counterMetrics, storage.CounterMetric{
			ID:    tmp.(storage.CounterMetric).ID,
			Delta: tmp.(storage.CounterMetric).Delta,
		})
	}

	return &gaugeMetrics, &counterMetrics, nil
}
