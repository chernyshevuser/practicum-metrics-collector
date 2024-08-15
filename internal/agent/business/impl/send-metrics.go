package impl

import (
	"context"
	"net/http"
	"time"
)

func (s *svc) sendMetrics(ctx context.Context) {
	metrics := func() []Metric {
		s.mu.Lock()
		defer s.mu.Unlock()

		res := make([]Metric, len(s.metrics))
		copy(res, s.metrics)

		return res
	}()

	cl := http.Client{}
	defer cl.CloseIdleConnections()

	for _, m := range metrics {
		err := s.sendWithRetry(ctx, &cl, m)
		if err != nil {
			s.logger.Errorw(
				"can't send metric",
				"name", m.ID,
				"type", m.Type,
				"val", m.Val,
				"reason", err,
			)
		}
	}
}

type SendReq struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func (s *svc) sendWithRetry(ctx context.Context, cl *http.Client, m Metric) (err error) {
	timeouts := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}

	for attempt := 0; attempt < len(timeouts); attempt++ {
		err = func() error {
			//TODO fix
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "addr", nil)
			if err != nil {
				s.logger.Errorw(
					"error in creating request",
					"reason", err.Error(),
				)
				return err
			}

			// if config.HashKey != "" {
			// 	sign := hash.Sign(reqByte, config.HashKey)
			// 	req.Header.Set("HashSHA256", base64.StdEncoding.EncodeToString(sign))
			// }

			req.Header.Set("Content-Type", "application/json")
			// req.Header.Set("Content-Encoding", "gzip")
			// req.Header.Set("Accept-Encoding", "gzip")

			resp, err := cl.Do(req)
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
		}()
		if err != nil {
			s.logger.Errorw(
				"error in sending",
				"reason:", err.Error(),
				"sleep:", timeouts[attempt],
			)

			time.Sleep(timeouts[attempt])
		} else {
			return
		}
	}

	if err != nil {
		return err
	}

	return nil
}
