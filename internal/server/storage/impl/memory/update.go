package memorystorage

import (
	"context"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) Update(ctx context.Context, gaugeMetrics []storage.GaugeMetric, counterMetrics []storage.CounterMetric) error {
	for _, m := range gaugeMetrics {
		s.gaugeStorage.Set(m.ID, m)
	}

	for _, m := range counterMetrics {
		storedVal := decimal.NewFromInt(0)
		stored, ok := s.counterStorage.Get(m.ID)
		if ok {
			storedVal = stored.(decimal.Decimal)
		}

		s.counterStorage.Set(m.ID, storage.CounterMetric{
			ID:    m.ID,
			Delta: storedVal.Add(m.Delta),
		})
	}

	return nil
}
