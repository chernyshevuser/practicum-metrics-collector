package db

import (
	"context"

	files "github.com/chernyshevuser/practicum-metrics-collector"
	"github.com/pressly/goose/v3"

	"github.com/chernyshevuser/practicum-metrics-collector/tools/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func (s *svc) Actualize(ctx context.Context) error {
	return s.conn.AcquireFunc(ctx, func(*pgxpool.Conn) error {
		goose.SetBaseFS(files.Migrations)

		if err := goose.SetDialect("pgx"); err != nil {
			return err
		}

		con, err := goose.OpenDBWithDriver(
			"pgx",
			config.DatabaseURI,
		)
		if err != nil {
			return err
		}

		if err := goose.Up(con, "migrations"); err != nil {
			return err
		}

		if err := con.Close(); err != nil {
			return err
		}

		return nil
	})
}
