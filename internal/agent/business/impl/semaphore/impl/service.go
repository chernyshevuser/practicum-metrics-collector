package impl

import "github.com/chernyshevuser/practicum-metrics-collector/internal/agent/business/impl/semaphore"

type service struct {
	semaCh chan struct{}
}

func (s *service) Acquire() {
	s.semaCh <- struct{}{}
}

func (s *service) Release() {
	<-s.semaCh
}

func (s *service) Close() {
	close(s.semaCh)
}

func New(maxReq int64) semaphore.Semaphore {
	return &service{
		semaCh: make(chan struct{}, maxReq),
	}
}
