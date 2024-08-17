package memorystorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
	"github.com/shopspring/decimal"
)

func (s *svc) Lock() {
	s.mu.Lock()
}

func (s *svc) Unlock() {
	s.mu.Unlock()
}

func (s *svc) Ping(ctx context.Context) error {
	return nil
}

type rawData struct {
	Metrics []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		Val  string `json:"val"`
	}
}

func (s *svc) Actualize(ctx context.Context) error {
	file, err := os.OpenFile(s.filepath, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			s.logger.Infow(
				"db is not actualized",
				"source file", s.filepath,
				"reason", "file doesn't exist",
			)
			return nil
		}

		return err
	}
	defer file.Close()

	encoder := json.NewDecoder(file)

	var rawData rawData

	if err = encoder.Decode(&rawData); err != nil {
		return err
	}

	metrics := make([]storage.Metric, 0, len(rawData.Metrics))
	for _, m := range rawData.Metrics {
		val, err := decimal.NewFromString(m.Val)
		if err != nil {
			return fmt.Errorf("can't parse string to decimal, reason: %v", err)
		}

		metrics = append(
			metrics,
			storage.Metric{
				ID:   m.ID,
				Type: m.Type,
				Val:  val,
			},
		)
	}

	err = s.Set(ctx, metrics)
	if err != nil {
		return fmt.Errorf("can't set metrics to db, reason: %v", err)
	}

	return nil
}

func (s *svc) Dump(ctx context.Context) error {
	metrics, err := s.GetAll(ctx)
	if err != nil {
		return fmt.Errorf("can't get all metrics from db, reason: %v", err)
	}

	//no metrics to dump
	if metrics == nil {
		return nil
	}

	var rawData rawData
	for _, m := range *metrics {
		rawData.Metrics = append(
			rawData.Metrics, struct {
				ID   string `json:"id"`
				Type string `json:"type"`
				Val  string `json:"val"`
			}{
				ID:   m.ID,
				Type: m.Type,
				Val:  m.Val.String(),
			},
		)
	}

	file, err := os.OpenFile(s.filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(rawData); err != nil {
		return fmt.Errorf("can't encode metrics, reason: %v", err)
	}

	s.logger.Infow(fmt.Sprintf("all metrics were dumpled to %s", s.filepath))

	return nil
}

func (s *svc) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.filepath != "" {
		err := s.Dump(context.Background())
		if err != nil {
			return err
		}
	}

	s.logger.Info("goodbye from db-svc")

	return nil
}
