package impl

import (
	"context"
	"time"
)

func (s *svc) Run(ctx context.Context) {
	updateDefaultTicker := time.NewTicker(time.Duration(s.updateInterval) * time.Second)
	updateExtraTicker := time.NewTicker(time.Duration(s.updateInterval) * time.Second)
	sendTicker := time.NewTicker(time.Duration(s.sendInterval) * time.Second)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				s.logger.Infow("close done")
				return
			case <-ctx.Done():
				s.logger.Infow("ctx done")
				return
			case <-updateDefaultTicker.C:
				s.logger.Infow(
					"update metrics",
					"status", "start",
				)
				s.collectMetrics()
				s.logger.Infow(
					"update metrics",
					"status", "finished",
				)
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				s.logger.Infow("close done")
				return
			case <-ctx.Done():
				s.logger.Infow("ctx done")
				return
			case <-updateExtraTicker.C:
				s.logger.Infow(
					"update extra metrics",
					"status", "start",
				)
				s.collectExtraMetrics()
				s.logger.Infow(
					"update extra metrics",
					"status", "finished",
				)
			}
		}
	}()

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			select {
			case <-s.closeCh:
				s.logger.Infow("close done")
				return
			case <-ctx.Done():
				s.logger.Infow("ctx done")
				return
			case <-sendTicker.C:
				s.logger.Infow(
					"send metrics",
					"status", "start",
				)
				s.wg.Add(1)
				go func() {
					defer s.wg.Done()

					s.semaphore.Acquire()
					defer s.semaphore.Release()

					s.sendMetrics(ctx)
					s.logger.Infow(
						"send metrics",
						"status", "finished",
					)
				}()
			}
		}
	}()

	s.wg.Wait()
}
