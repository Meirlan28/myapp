package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (r *UserRepository) GetByID(
	ctx context.Context,
	id int,
) (*domain.User, error) {
	query := `
		SELECT id, name, age
		FROM myapp.users
		WHERE id = $1
	`

	var u domain.User

	if err := r.db.QueryRow(ctx, query, id).Scan(
		&u.ID,
		&u.Name,
		&u.Age,
	); err != nil {
		return nil, postgres.MapError(err)
	}

	return &u, nil
}
