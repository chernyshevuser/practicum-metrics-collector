package db

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5"
)

func (s *svc) Lock() {
	s.mu.Lock()
}

func (s *svc) Unlock() {
	s.mu.Unlock()
}

func (s *svc) Ping(ctx context.Context) error {
	return s.wrap(func() error {
		return s.conn.Ping(ctx)
	})
}

func (s *svc) Dump(ctx context.Context) error {
	return nil
}

func (s *svc) Close() error {
	s.conn.Close()

	s.logger.Info("goodbye from db-svc")

	return nil
}

func (s *svc) beginR(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadOnly,
	})
	return tx, err
}

func (s *svc) beginW(ctx context.Context) (pgx.Tx, error) {
	tx, err := s.conn.BeginTx(ctx, pgx.TxOptions{
		AccessMode: pgx.ReadWrite,
	})
	return tx, err
}

func (s *svc) wrap(f func() error) error {
	var (
		timeouts = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
		err      error
		pgErr    *pgconn.PgError
	)

	for i := 0; i < len(timeouts); i++ {
		err = f()
		if err == nil {
			return nil
		}

		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			time.Sleep(timeouts[i])
			continue
		}

		return err
	}

	return err
}
