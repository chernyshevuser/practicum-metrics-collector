package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/shopspring/decimal"
)

func (s *svc) GetMetricValue(ctx context.Context, metricType, metricName string) (*decimal.Decimal, error) {
	t := s.parseMetricType(metricType)

	s.logger.Infow(
		"parse metric",
		"given type", metricType,
		"parsed type", t,
		"given name", metricName,
	)

	//case unknown
	if t == business.Unknown {
		return nil, fmt.Errorf("given metric type(%s) in unknown", metricType)
	}

	//case counter
	if t == business.Counter {
		m, err := s.db.GetCounter(ctx, metricName)
		if err != nil {
			s.logger.Errorw(
				"storage problem",
				"msg", "can't get counter metric val",
				"reason", err,
			)
			return nil, fmt.Errorf("can't get counter metric val from db, reason: %v", err)
		}

		if m == nil {
			return nil, nil
		}

		return &m.Delta, nil
	}

	//case gauge
	m, err := s.db.GetGauge(ctx, metricName)
	if err != nil {
		s.logger.Errorw(
			"storage problem",
			"msg", "can't get gauge metric val",
			"reason", err,
		)
		return nil, fmt.Errorf("can't get gauge metric val from db, reason: %v", err)
	}

	if m == nil {
		return nil, nil
	}

	return &m.Value, nil
}

func (s *svc) GetMetricValueJSON() error {
	panic("")
}

func (s *svc) GetAllMetrics(ctx context.Context) ([]business.CounterMetric, []business.GaugeMetric, error) {
	gaugeMetrics, counterMetrics, err := s.db.GetAll(ctx)
	if err != nil {
		s.logger.Errorw(
			"storage problem",
			"msg", "can't get all metrics vals",
			"reason", err,
		)
		return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't get metrics from db, reason: %v", err)
	}

	gaugeMetricsRes := make([]business.GaugeMetric, 0, len(*gaugeMetrics))
	for _, tmp := range *gaugeMetrics {
		gaugeMetricsRes = append(gaugeMetricsRes, business.GaugeMetric{
			ID:    tmp.ID,
			Value: tmp.Value,
		})
	}

	counterMetricsRes := make([]business.CounterMetric, 0, len(*counterMetrics))
	for _, tmp := range *counterMetrics {
		counterMetricsRes = append(counterMetricsRes, business.CounterMetric{
			ID:    tmp.ID,
			Delta: tmp.Delta,
		})
	}

	return counterMetricsRes, gaugeMetricsRes, nil
}
