package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sso/pkg/config"
)

type Database struct {
	cluster *pgxpool.Pool
	logger  *slog.Logger
}

func NewDatabase(ctx context.Context, cfg *config.Config, logger *slog.Logger) *Database {
	pool, err := pgxpool.New(ctx, generateDsn(cfg))
	if err != nil {
		panic(err)
	}

	return &Database{
		cluster: pool,
		logger:  logger,
	}
}

func generateDsn(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.DBName)
}

// Read одного объекта в структуру.
func (db *Database) Get(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	const operation = "Postgres.Get"
	logger := db.logger.With(slog.String("operation", operation))

	logger.Info("Execute query", "query", query, "args", args)
	err := pgxscan.Get(ctx, db.cluster, dst, query, args...)
	if err != nil {
		logger.Error("Failed to execute query", "query", query, "error", err)
	}

	return err
}

// Read нескольких объектов одного типа в срез структур.
// TODO
// Проблема с args. Если их нет, то выплывает ошибка:
// [] scany: query multiple result rows: expected 0 arguments, got 1.
// Сейчас сделано, на мой взгляд, коряво.
func (db *Database) Select(ctx context.Context, dst interface{}, query string, args ...interface{}) error {
	const operation = "Postgres.Select"
	logger := db.logger.With(slog.String("operation", operation))

	logger.Info("Execute query", "query", query, "args", args)
	var err error
	if args == nil {
		err = pgxscan.Select(ctx, db.cluster, dst, query)
	} else {
		err = pgxscan.Select(ctx, db.cluster, dst, query, args)
	}

	if err != nil {
		logger.Error("Failed to execute query", "query", query, "error", err)
	}

	return err
}

// Create, Update, Delete.
func (db *Database) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	const operation = "Postgres.Exec"
	logger := db.logger.With(slog.String("operation", operation))

	logger.Info("Execute query", "query", query, "args", args)
	ct, err := db.cluster.Exec(ctx, query, args)
	if err != nil {
		logger.Error("Failed to execute query", "query", query, "error", err)
	}

	return ct, err
}

// Read нескольких РАЗНЫХ параметров.
func (db *Database) ExecQueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	const operation = "Postgres.ExecQueryRow"
	logger := db.logger.With(slog.String("operation", operation))

	// TODO обработать ошибку
	logger.Info("Execute query", "query", query, "args", args)
	row := db.cluster.QueryRow(ctx, query, args...)
	/*if errors.Is(, pgx.ErrNoRows) {
		logger.Error("Failed to execute query", "query", query, "error", err)
	}*/

	return row
}
