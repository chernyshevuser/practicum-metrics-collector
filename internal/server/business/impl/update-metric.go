package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) UpdateMetrics(ctx context.Context, metrics []business.RawMetric) (updatedCounterMetrics []business.CounterMetric, updatedGaugeMetrics []business.GaugeMetric, err error) {
	var counterMetrics []storage.CounterMetric
	var gaugeMetrics []storage.GaugeMetric

	for _, metric := range metrics {
		t := s.parseMetricType(metric.Type)

		//case unknown type
		if t == business.Unknown {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("given metric type(%s) in unknown", metric.Type)
		}

		if t == business.Counter {
			val, err := decimal.NewFromString(metric.Value)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse counter metric value(%s) to decimal.Decimal, reason: %v", metric.Value, err)
			}

			if !s.isDecimalInt(val) {
				return []business.CounterMetric{}, []business.GaugeMetric{}, business.ErrWrongMetricVal
			}

			counterMetrics = append(
				counterMetrics,
				storage.CounterMetric{
					ID:    metric.ID,
					Delta: val,
				},
			)
		} else if t == business.Gauge {
			val, err := decimal.NewFromString(metric.Value)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse gauge metric value(%s) to decimal.Decimal, reason: %v", metric.Value, err)
			}

			gaugeMetrics = append(
				gaugeMetrics,
				storage.GaugeMetric{
					ID:    metric.ID,
					Value: val,
				},
			)
		} else {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("smth wrong with parsing metric type: %s", metric.Type)
		}
	}

	s.db.Lock()
	defer s.db.Unlock()

	for i := range counterMetrics {
		stored, err := s.db.GetCounter(ctx, counterMetrics[i].ID)
		if err != nil {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't get metric from db, reason: %v", err)
		}

		if stored != nil {
			counterMetrics[i].Delta = counterMetrics[i].Delta.Add(stored.Delta)
		}
	}

	err = s.db.Update(ctx, gaugeMetrics, counterMetrics)
	if err != nil {
		return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't update metrics, reason: %v", err)
	}

	return func() []business.CounterMetric {
			res := make([]business.CounterMetric, 0, len(counterMetrics))
			for _, m := range counterMetrics {
				res = append(res, business.CounterMetric{
					ID:    m.ID,
					Delta: m.Delta,
				})
			}
			return res
		}(),
		func() []business.GaugeMetric {
			res := make([]business.GaugeMetric, 0, len(gaugeMetrics))
			for _, m := range gaugeMetrics {
				res = append(res, business.GaugeMetric{
					ID:    m.ID,
					Value: m.Value,
				})
			}
			return res
		}(),
		nil
}
