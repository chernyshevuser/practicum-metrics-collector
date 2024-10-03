package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) GetMetricValue(ctx context.Context, metricType, metricName string) (val *decimal.Decimal, mType business.MetricType, err error) {
	s.db.Lock()
	defer s.db.Unlock()

	t := parseMetricType(metricType)

	//case unknown
	if t == business.Unknown {
		return nil, "", fmt.Errorf("given metric type(%s) in unknown", metricType)
	}

	//case counter
	if t == business.Counter {
		var stored *storage.Metric
		stored, err = s.db.Get(ctx, storage.BuildKey(metricName, string(t)))
		if err != nil {
			s.logger.Errorw(
				"storage problem",
				"msg", "can't get counter metric val",
				"reason", err,
			)
			return nil, "", fmt.Errorf("can't get counter metric val from db, reason: %v", err)
		}

		if stored == nil {
			return nil, "", nil
		}

		delta := decimal.NewFromInt(stored.Delta)
		return &delta, business.Counter, nil
	}

	//case gauge
	stored, err := s.db.Get(ctx, storage.BuildKey(metricName, string(t)))
	if err != nil {
		s.logger.Errorw(
			"storage problem",
			"msg", "can't get gauge metric val",
			"reason", err,
		)
		return nil, "", fmt.Errorf("can't get gauge metric val from db, reason: %v", err)
	}

	if stored == nil {
		return nil, "", nil
	}

	value := decimal.NewFromFloat(stored.Val)
	return &value, business.Gauge, nil
}

func (s *svc) GetAllMetrics(ctx context.Context) ([]business.CounterMetric, []business.GaugeMetric, error) {
	s.db.Lock()
	defer s.db.Unlock()

	metrics, err := s.db.GetAll(ctx)
	if err != nil {
		s.logger.Errorw(
			"storage problem",
			"msg", "can't get all metrics vals",
			"reason", err,
		)
		return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't get metrics from db, reason: %v", err)
	}

	var (
		counterMetrics []business.CounterMetric
		gaugeMetrics   []business.GaugeMetric
	)

	for _, tmp := range *metrics {
		if tmp.Type == string(business.Counter) {
			counterMetrics = append(counterMetrics, business.CounterMetric{
				ID:    tmp.ID,
				Delta: decimal.NewFromInt(tmp.Delta),
			})
		} else if tmp.Type == string(business.Gauge) {
			gaugeMetrics = append(gaugeMetrics, business.GaugeMetric{
				ID:    tmp.ID,
				Value: decimal.NewFromFloat(tmp.Val),
			})
		} else {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("incorrect metric type(%s) from db", tmp.Type)
		}
	}

	return counterMetrics, gaugeMetrics, nil
}
