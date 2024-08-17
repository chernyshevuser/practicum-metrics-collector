package storage

import (
	"context"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Metric struct {
	ID   string
	Type string
	Val  decimal.Decimal
}

func BuildKey(metricID, metricType string) string {
	return fmt.Sprintf("%s_%s", metricID, metricType)
}

func ParseKey(key string) (metricID, metricType string, err error) {
	tmp := strings.Split(key, "_")
	if len(tmp) != 2 {
		return "", "", fmt.Errorf("can't parse key")
	}

	return tmp[0], tmp[1], nil
}

type Storage interface {
	Set(ctx context.Context, metrics []Metric) (err error)

	Get(ctx context.Context, key string) (*Metric, error)
	GetAll(ctx context.Context) (*[]Metric, error)

	Lock()
	Unlock()

	Actualize(ctx context.Context) error
	Dump(ctx context.Context) error

	Ping(ctx context.Context) error
	Close() error
}
