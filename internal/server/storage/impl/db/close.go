package db

func (s *svc) Close() error {
	s.conn.Close()
	s.logger.Info("goodbye from db-svc")
	return nil
}
