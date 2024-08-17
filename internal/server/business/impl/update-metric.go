package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) UpdateMetrics(ctx context.Context, rawMetrics []business.RawMetric) (updatedCounterMetrics []business.CounterMetric, updatedGaugeMetrics []business.GaugeMetric, err error) {
	var counterMetrics, gaugeMetrics []storage.Metric

	for _, rm := range rawMetrics {
		t := s.parseMetricType(rm.Type)

		//case unknown type
		if t == business.Unknown {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("given metric type(%s) in unknown", rm.Type)
		}

		if t == business.Counter {
			val, err := decimal.NewFromString(rm.Value)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse counter metric value(%s) to decimal.Decimal, reason: %v", rm.Value, err)
			}

			if !s.isDecimalInt(val) {
				return []business.CounterMetric{}, []business.GaugeMetric{}, business.ErrWrongMetricVal
			}

			counterMetrics = append(
				counterMetrics, storage.Metric{
					ID:   rm.ID,
					Type: rm.Type,
					Val:  val,
				},
			)
		} else if t == business.Gauge {
			val, err := decimal.NewFromString(rm.Value)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse gauge metric value(%s) to decimal.Decimal, reason: %v", rm.Value, err)
			}

			gaugeMetrics = append(
				gaugeMetrics, storage.Metric{
					ID:   rm.ID,
					Type: rm.Type,
					Val:  val,
				},
			)
		} else {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("smth wrong with parsing metric type: %s", rm.Type)
		}
	}

	s.db.Lock()
	defer s.db.Unlock()

	for i := range counterMetrics {
		stored, err := s.db.Get(ctx, storage.BuildKey(counterMetrics[i].ID, counterMetrics[i].Type))
		if err != nil {
			s.logger.Errorw(
				"can't get counter metric from db",
				"reason", err,
			)
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't get metric from db, reason: %v", err)
		}

		if stored != nil {
			counterMetrics[i].Val = counterMetrics[i].Val.Add(stored.Val)
		}
	}

	err = s.db.Set(ctx, counterMetrics)
	if err != nil {
		s.logger.Errorw(
			"can't update counter metrics",
			"reason", err,
		)
		return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't update counter metrics, reason: %v", err)
	}

	err = s.db.Set(ctx, gaugeMetrics)
	if err != nil {
		s.logger.Errorw(
			"can't update gauge metrics",
			"reason", err,
		)
		return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't update gauge metrics, reason: %v", err)
	}

	return func() []business.CounterMetric {
			res := make([]business.CounterMetric, 0, len(counterMetrics))
			for _, m := range counterMetrics {
				res = append(res, business.CounterMetric{
					ID:    m.ID,
					Delta: m.Val,
				})
			}
			return res
		}(),
		func() []business.GaugeMetric {
			res := make([]business.GaugeMetric, 0, len(gaugeMetrics))
			for _, u := range gaugeMetrics {
				res = append(res, business.GaugeMetric{
					ID:    u.ID,
					Value: u.Val,
				})
			}
			return res
		}(),
		nil
}
