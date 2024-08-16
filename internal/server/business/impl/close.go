package impl

func (s *svc) Close() {
	s.logger.Info("goodbye from business-svc")
}
