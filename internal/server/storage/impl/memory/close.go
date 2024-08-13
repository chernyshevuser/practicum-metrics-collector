package memorystorage

func (s *svc) Close() error {
	s.logger.Info("goodbye from db-svc")

	return nil
}
