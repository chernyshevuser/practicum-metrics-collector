package impl

import (
	"testing"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/business"
	"github.com/test-go/testify/assert"
)

func TestParseMetricType(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		expected business.MetricType
		_        string
	}{
		{
			name:     "valid counter",
			in:       "counter",
			expected: business.Counter,
		},
		{
			name:     "valid gauge",
			in:       "gauge",
			expected: business.Gauge,
		},
		{
			name:     "invalid",
			in:       "invalid",
			expected: business.Unknown,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			parsed := parseMetricType(test.in)
			assert.Equal(t, parsed, test.expected)
		})
	}
}
