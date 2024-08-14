package impl

import (
	"context"
	"fmt"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) UpdateMetric(ctx context.Context, metricType, metricName, metricVal string) error {
	t := s.parseMetricType(metricType)

	//case unknown
	if t == business.Unknown {
		return fmt.Errorf("given metric type(%s) in unknown", metricType)
	}

	//case counter
	if t == business.Counter {
		val, err := decimal.NewFromString(metricVal)
		if err != nil {
			return fmt.Errorf("can't parse counter metric value(%s) to decimal.Decimal, reason: %v", metricVal, err)
		}

		if !s.isDecimalInt(val) {
			return business.ErrWrongMetricVal
		}

		err = s.db.UpdateCounter(
			ctx,
			storage.CounterMetric{
				ID:    metricName,
				Delta: val,
			},
		)
		if err != nil {
			return fmt.Errorf("can't update counter metric in db, reason: %v", err)
		}

		return nil
	}

	//case gauge
	val, err := decimal.NewFromString(metricVal)
	if err != nil {
		return fmt.Errorf("can't parse gauge metric value(%s) to decimal.Decimal, reason: %v", metricVal, err)
	}

	err = s.db.UpdateGauge(
		ctx,
		storage.GaugeMetric{
			ID:    metricName,
			Value: val,
		},
	)
	if err != nil {
		return fmt.Errorf("can't update gauge metric in db, reason: %v", err)
	}

	return nil
}

func (s *svc) UpdateMetricJSON() error {
	panic("")
}

func (s *svc) UpdateMetricsJSON() error {
	panic("")
}
