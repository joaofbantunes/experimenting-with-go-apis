package internal

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

func CreatePool(ctx context.Context, config *Config) (*pgxpool.Pool, error) {
	replacer := strings.NewReplacer(
		"{user}", config.Database.User,
		"{password}", config.Database.Password,
	)
	pool, err := pgxpool.New(ctx, replacer.Replace(config.Database.BaseConnStr))
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func MigrateDB(ctx context.Context, pool *pgxpool.Pool, logger *slog.Logger) error {
	// go migrate expects a database/sql connection, so we need to get a *sql.DB from pgxpool.Pool
	conn := stdlib.OpenDBFromPool(pool)

	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			logger.ErrorContext(ctx, "Failed to close database connection", slog.Any("error", err))
		}
	}(conn)

	driver, err := pgx.WithInstance(conn, &pgx.Config{})
	if err != nil {
		logger.ErrorContext(ctx, "pgx.WithInstance", slog.Any("error", err))
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"experimenting_with_std_apis",
		driver)

	if err != nil {
		logger.ErrorContext(ctx, "migrate.NewWithDatabaseInstance", slog.Any("error", err))
		return err
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.ErrorContext(ctx, "migrate.Up", slog.Any("error", err))
		return err
	}

	return nil
}
