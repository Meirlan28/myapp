package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgres(
	ctx context.Context,
	postgresConfig Config,
) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(postgresConfig.DSN())
	if err != nil {
		return nil, fmt.Errorf("parse postgres config: %w", err)
	}

	cfg.MaxConns = postgresConfig.MaxConns

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()

		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return pool, nil
}
