package impl

import (
	"context"
	"fmt"
	"strconv"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) UpdateMetrics(ctx context.Context, rawMetrics []business.RawMetric) (updatedCounterMetrics []business.CounterMetric, updatedGaugeMetrics []business.GaugeMetric, err error) {
	var counterMetrics, gaugeMetrics []storage.Metric

	for _, rm := range rawMetrics {
		t := parseMetricType(rm.Type)

		//case unknown type
		if t == business.Unknown {
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("given metric type(%s) in unknown", rm.Type)
		}

		if t == business.Counter {
			var delta int64
			delta, err = strconv.ParseInt(rm.Value, 10, 64)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse counter metric value(%s) to int64, reason: %v", rm.Value, err)
			}

			counterMetrics = append(
				counterMetrics, storage.Metric{
					ID:    rm.ID,
					Type:  rm.Type,
					Delta: delta,
				},
			)
		} else if t == business.Gauge {
			var val float64
			val, err = strconv.ParseFloat(rm.Value, 64)
			if err != nil {
				return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't parse gauge metric value(%s) to float64, reason: %v", rm.Value, err)
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

	// make unique
	counterMetrics = func() []storage.Metric {
		data := make(map[string]int64)
		for _, cm := range counterMetrics {
			data[cm.ID] += cm.Delta
		}

		out := make([]storage.Metric, 0, len(data))
		for k, v := range data {
			out = append(out, storage.Metric{
				ID:    k,
				Type:  string(business.Counter),
				Delta: v,
			})
		}

		return out
	}()

	// make unique
	gaugeMetrics = func() []storage.Metric {
		data := make(map[string]float64)
		for _, cm := range gaugeMetrics {
			data[cm.ID] = cm.Val
		}

		out := make([]storage.Metric, 0, len(data))
		for k, v := range data {
			out = append(out, storage.Metric{
				ID:   k,
				Type: string(business.Gauge),
				Val:  v,
			})
		}

		return out
	}()

	s.db.Lock()
	defer s.db.Unlock()

	// increment counter vals
	for i := range counterMetrics {
		var stored *storage.Metric
		stored, err = s.db.Get(ctx, storage.BuildKey(counterMetrics[i].ID, counterMetrics[i].Type))
		if err != nil {
			s.logger.Errorw(
				"can't get counter metric from db",
				"reason", err,
			)
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't get metric from db, reason: %v", err)
		}

		if stored != nil {
			counterMetrics[i].Delta += stored.Delta
		}
	}

	for _, m := range counterMetrics {
		if err = s.db.Set(ctx, m); err != nil {
			s.logger.Errorw(
				"can't update counter metrics",
				"reason", err,
			)
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't update counter metrics, reason: %v", err)
		}
	}

	for _, m := range gaugeMetrics {
		if err = s.db.Set(ctx, m); err != nil {
			s.logger.Errorw(
				"can't update gauge metrics",
				"reason", err,
			)
			return []business.CounterMetric{}, []business.GaugeMetric{}, fmt.Errorf("can't update gauge metrics, reason: %v", err)
		}
	}

	return func() []business.CounterMetric {
			res := make([]business.CounterMetric, 0, len(counterMetrics))
			for _, m := range counterMetrics {
				res = append(res, business.CounterMetric{
					ID:    m.ID,
					Delta: decimal.NewFromInt(m.Delta),
				})
			}
			return res
		}(),
		func() []business.GaugeMetric {
			res := make([]business.GaugeMetric, 0, len(gaugeMetrics))
			for _, m := range gaugeMetrics {
				res = append(res, business.GaugeMetric{
					ID:    m.ID,
					Value: decimal.NewFromFloat(m.Val),
				})
			}
			return res
		}(),
		nil
}
