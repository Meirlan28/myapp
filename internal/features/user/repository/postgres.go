package repository

import (
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db     *pgxpool.Pool
	logger *slog.Logger
}

func New(db *pgxpool.Pool, logger *slog.Logger) *UserRepository {
	return &UserRepository{db, logger}
}
