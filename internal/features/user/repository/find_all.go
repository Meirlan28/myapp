package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
)

func (ur *UserRepository) FindAll(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, error) {
	ur.logger.Info("FindAll", "min_age", minAge, "max_age", maxAge, "limit", limit, "offset", offset)
	query := `
	SELECT users.id, users.name, users.age 
	FROM myapp.users 
	WHERE ($1::int IS NULL OR users.age >= $1)
  	AND ($2::int IS NULL OR users.age <= $2)
	LIMIT $3
	OFFSET $4
	`
	rows, err := ur.db.Query(ctx, query, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, DatabaseError
	}
	defer rows.Close()

	var users []domain.User

	for rows.Next() {
		var u domain.User

		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Age,
		); err != nil {
			return nil, DatabaseError
		}

		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, DatabaseError
	}

	return users, nil
}
