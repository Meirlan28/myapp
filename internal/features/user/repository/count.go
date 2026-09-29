package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) Count(
	ctx context.Context,
	minAge *int,
	maxAge *int,
) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM myapp.users
		WHERE ($1::int IS NULL OR age >= $1)
		  AND ($2::int IS NULL OR age <= $2)
	`

	var count int
	if err := r.db.QueryRow(ctx, query, minAge, maxAge).Scan(&count); err != nil {
		return 0, postgres.MapError(err)
	}

	return count, nil
}
