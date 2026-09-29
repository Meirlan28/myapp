package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) Update(
	ctx context.Context,
	id int,
	u *domain.User,
) (*domain.User, error) {
	query := `
		UPDATE myapp.users
		SET name = $1, age = $2
		WHERE id = $3
		RETURNING id, name, age
	`

	var updated domain.User

	if err := r.db.QueryRow(
		ctx,
		query,
		u.Name,
		u.Age,
		id,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Age,
	); err != nil {
		return nil, postgres.MapError(err)
	}

	return &updated, nil
}
