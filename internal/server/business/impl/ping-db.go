package impl

import (
	"context"
	"fmt"
)

func (s *svc) PingDB(ctx context.Context) error {
	err := s.db.Ping(ctx)
	if err != nil {
		s.logger.Errorw(
			"can't ping db",
			"reason", err,
		)
		return fmt.Errorf("can't ping db")
	}

	return nil
}
