package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) List(
	ctx context.Context,
	minAge *int,
	maxAge *int,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	query := `
		SELECT id, name, age
		FROM myapp.users
		WHERE ($1::int IS NULL OR age >= $1)
		  AND ($2::int IS NULL OR age <= $2)
		LIMIT $3
		OFFSET $4
	`

	rows, err := r.db.Query(ctx, query, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, postgres.MapError(err)
	}
	defer rows.Close()

	users := make([]domain.User, 0)

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Age,
		); err != nil {
			return nil, postgres.MapError(err)
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.MapError(err)
	}

	return users, nil
}
