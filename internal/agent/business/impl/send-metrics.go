package impl

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/config"
	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/utils"
	"github.com/chernyshevuser/practicum-metrics-collector/tools/crypto"
)

func (s *svc) sendMetrics(ctx context.Context) {
	metrics := func() []Metric {
		s.mu.Lock()
		defer s.mu.Unlock()

		var res []Metric

		res = append(res, s.metrics...)
		res = append(res, s.extraMetrics...)

		return res
	}()

	cl := http.Client{}
	defer cl.CloseIdleConnections()

	err := s.sendWithRetry(ctx, &cl, metrics)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}

		s.logger.Errorw(
			"can't send metrics",
			"reason", err,
		)
	}
}

type unit struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type sendReq []unit

func (s *svc) sendWithRetry(ctx context.Context, cl *http.Client, metrics []Metric) error {
	var req sendReq

	for _, m := range metrics {
		u := unit{
			ID:    m.ID,
			MType: m.Type,
		}

		if m.Type == string(business.CounterMT) {
			valInt64 := m.Val.IntPart()
			u.Delta = &valInt64
		} else if m.Type == string(business.GaugeMT) {
			valFloat64 := m.Val.InexactFloat64()
			u.Value = &valFloat64
		} else {
			s.logger.Errorw(
				"wrong metric type",
				"type", m.Type,
				"ID", m.ID,
			)
			return fmt.Errorf("wrong metric type")
		}

		req = append(req, u)
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return err
	}

	buf, err := utils.Compress(reqByte)
	if err != nil {
		return err
	}

	timeouts := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for attempt := 0; attempt < len(timeouts); attempt++ {
		if err = func() error {
			var req *http.Request
			req, err = http.NewRequestWithContext(ctx, http.MethodPost, s.addr, buf)
			if err != nil {
				s.logger.Errorw(
					"error in creating request",
					"reason", err.Error(),
				)
				return err
			}

			if config.HashKey != "" {
				sign := crypto.Sign(reqByte, config.HashKey)
				req.Header.Set("HashSHA256", base64.StdEncoding.EncodeToString(sign))
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Content-Encoding", "gzip")
			req.Header.Set("Accept-Encoding", "gzip")

			var resp *http.Response
			resp, err = cl.Do(req)
			if err != nil {
				s.logger.Errorw(
					"error in sending request",
					"reason", err.Error(),
				)
				return err
			}

			if resp != nil {
				s.logger.Infow(
					"response",
					"status", resp.StatusCode,
				)

				if err = resp.Body.Close(); err != nil {
					return err
				}
			}
			return err
		}(); err == nil {
			return nil
		}

		if errors.Is(err, context.Canceled) {
			return err
		}

		s.logger.Errorw(
			"error in sending",
			"reason:", err.Error(),
			"sleep:", timeouts[attempt],
		)

		time.Sleep(timeouts[attempt])
	}
	if err != nil {
		return err
	}

	return nil
}
