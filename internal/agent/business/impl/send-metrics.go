package impl

import (
	"context"
	"net/http"
	"sync"

	"github.com/chernyshevuser/practicum-metrics-collector/internal/agent/constants"
)

// TODO not finished
func (s *svc) SendMetrics(ctx context.Context) {
	s.mu.Lock()

	metrics := make([]Metric, len(s.metrics))
	copy(metrics, s.metrics)

	s.mu.Unlock()

	ch := make(chan Metric, constants.Sz)
	wg := sync.WaitGroup{}
	cl := http.Client{}

	for i := 0; i < constants.Sz; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for m := range ch {
				err := sendWithRetry(&cl, m)
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
		}()
	}

	for _, m := range metrics {
		ch <- m
	}

	wg.Wait()
}

func sendWithRetry(cl *http.Client, m Metric) error {
	panic("")
}
