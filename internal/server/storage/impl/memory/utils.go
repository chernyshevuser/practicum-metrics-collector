package memorystorage

import "context"

func (s *svc) Lock() {
	s.mu.Lock()
}

func (s *svc) Unlock() {
	s.mu.Unlock()
}

func (s *svc) Ping(ctx context.Context) error {
	return nil
}

func (s *svc) Close() error {
	s.logger.Info("goodbye from db-svc")

	return nil
}
