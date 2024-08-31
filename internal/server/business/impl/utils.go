package impl

import (
	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
)

func (s *svc) parseMetricType(in string) business.MetricType {
	if in == string(business.Counter) {
		return business.Counter
	}

	if in == string(business.Gauge) {
		return business.Gauge
	}

	return business.Unknown
}
