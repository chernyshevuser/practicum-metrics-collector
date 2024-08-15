package impl

func (s *svc) Close() {
	s.closeCh <- struct{}{}
	close(s.closeCh)

	s.wg.Wait()

	s.logger.Info("goodbye from agent-svc")
}
