package db

import (
	"context"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jackc/pgx/v5"
)

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
	s.conn.Close()
	s.logger.Info("goodbye from db-svc")
	return nil
}

func (s *svc) query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := s.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (s *svc) queryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return s.conn.QueryRow(ctx, query, args...)
}

func (s *svc) exec(ctx context.Context, query string, args ...interface{}) error {
	_, err := s.conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}
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
