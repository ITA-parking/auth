package repo

import (
	"auth-service/util"
	"context"
	"database/sql"
	"errors"
	sq "github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

var psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
var Pool *pgxpool.Pool

func InitDB(dsn string) error {
	var err error
	Pool, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return errors.Join(errors.New("failed to connect to DB"), err)
	}

	err = Pool.Ping(context.Background())
	if err != nil {
		return errors.Join(errors.New("failed to ping DB"), err)
	}

	slog.Info("Connected to PostgreSQL")
	return migrateDB(dsn)
}

func migrateDB(dsn string) error {
	slog.Info("Starting migrations")
	conn, sqlErr := sql.Open("postgres", dsn)
	if sqlErr != nil {
		return errors.Join(errors.New("failed to connect to DB to apply migrations"), sqlErr)
	}

	driver, driverErr := postgres.WithInstance(conn, &postgres.Config{})
	if driverErr != nil {
		return errors.Join(errors.New("failed to create driver for DB to apply migrations"), driverErr)
	}

	migrationsPath, pathErr := util.GetAbsolutePath("./migrations")
	if pathErr != nil {
		return pathErr
	}

	m, migrationErr := migrate.NewWithDatabaseInstance(
		"file://"+migrationsPath,
		"postgres",
		driver,
	)
	if migrationErr != nil {
		return errors.Join(errors.New("failed to connect to DB to apply migrations"), migrationErr)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return errors.Join(errors.New("failed to apply migrations"), err)
	}

	return nil
}
