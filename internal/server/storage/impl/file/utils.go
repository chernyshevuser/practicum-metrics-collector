package filestorage

import "context"

func (s *svc) Lock() {
	panic("")
}

func (s *svc) Unlock() {
	panic("")
}

func (s *svc) Ping(ctx context.Context) error {
	return nil
}

func (s *svc) Close() error {
	panic("")
}
