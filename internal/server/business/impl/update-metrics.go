package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) UpdateMetrics(ctx context.Context, metrics []business.RawMetric) error {
	var counterMetrics []storage.CounterMetric
	var gaugeMetrics []storage.GaugeMetric

	for _, metric := range metrics {
		t := s.parseMetricType(metric.Type)

		//case unknown type
		if t == business.Unknown {
			return fmt.Errorf("given metric type(%s) in unknown", metric.Type)
		}

		if t == business.Counter {
			val, err := decimal.NewFromString(metric.Value)
			if err != nil {
				return fmt.Errorf("can't parse counter metric value(%s) to decimal.Decimal, reason: %v", metric.Value, err)
			}

			if !s.isDecimalInt(val) {
				return business.ErrWrongMetricVal
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
				return fmt.Errorf("can't parse gauge metric value(%s) to decimal.Decimal, reason: %v", metric.Value, err)
			}

			gaugeMetrics = append(
				gaugeMetrics,
				storage.GaugeMetric{
					ID:    metric.ID,
					Value: val,
				},
			)
		} else {
			return fmt.Errorf("smth wrong with parsing metric type: %s", metric.Type)
		}
	}

	err := s.db.Update(ctx, gaugeMetrics, counterMetrics)
	if err != nil {
		return fmt.Errorf("can't update metrics, reason: %v", err)
	}

	return nil
}
