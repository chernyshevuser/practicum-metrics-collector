package memorystorage

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/server/storage"
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
	ID    string  `json:"id"`
	Type  string  `json:"type"`
	Val   float64 `json:"val"`
	Delta int64   `json:"delta"`
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

	var rd []rawData

	if err = encoder.Decode(&rd); err != nil {
		return err
	}

	for _, m := range rd {
		err := s.Set(ctx, storage.Metric{
			ID:    m.ID,
			Type:  m.Type,
			Val:   m.Val,
			Delta: m.Delta,
		})
		if err != nil {
			return fmt.Errorf("can't set metrics to db, reason: %v", err)
		}
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

	rd := make([]rawData, 0, len(*metrics))

	for _, m := range *metrics {
		rd = append(
			rd, rawData{
				ID:    m.ID,
				Type:  m.Type,
				Val:   m.Val,
				Delta: m.Delta,
			},
		)
	}

	file, err := os.OpenFile(s.filepath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err = encoder.Encode(rd); err != nil {
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
