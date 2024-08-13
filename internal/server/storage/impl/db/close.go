package db

func (s *svc) Close() {
	s.conn.Close()

	s.logger.Info("goodbye from db-svc")
}
