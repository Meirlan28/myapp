package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(
	ctx context.Context,
	postgresConfig Config,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(postgresConfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	cfg.MaxConns = postgresConfig.MaxConns

	connectCtx, cancel := context.WithTimeout(ctx, postgresConfig.Timeout)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(connectCtx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(connectCtx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}

func MapError(err error) error {
	if err == nil {
		return nil
	}

	// Контекст запроса.
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		return errs.Unavailable("database timeout", err)

	case errors.Is(err, context.Canceled):
		return err
	}

	// Ничего не найдено.
	if errors.Is(err, pgx.ErrNoRows) {
		return errs.NotFound("resource not found", err)
	}

	// PostgreSQL-specific ошибки.
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505": // unique_violation
			return errs.Conflict("resource already exists", err)

		case "23503": // foreign_key_violation
			return errs.Invalid("related resource does not exist", err)

		case "23502": // not_null_violation
			return errs.Invalid("required field is missing", err)

		case "23514": // check_violation
			return errs.Invalid("invalid field value", err)

		case "40001": // serialization_failure
			return errs.Conflict("concurrent database operation conflict", err)

		case "40P01": // deadlock_detected
			return errs.Conflict("database operation conflict", err)

		case "53300", // too_many_connections
			"57P01", // admin_shutdown
			"57P02", // crash_shutdown
			"57P03": // cannot_connect_now
			return errs.Unavailable("database unavailable", err)
		}
	}

	return errs.Internal("database error", err)
}
