package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) Create(
	ctx context.Context,
	u *domain.User,
) (*domain.User, error) {
	query := `
		INSERT INTO myapp.users (name, age)
		VALUES ($1, $2)
		RETURNING id, name, age
	`

	if err := r.db.QueryRow(
		ctx,
		query,
		u.Name,
		u.Age,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Age,
	); err != nil {
		return nil, postgres.MapError(err)
	}

	return u, nil
}
