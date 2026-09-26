package repository

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (ur *UserRepository) FindByID(
	ctx context.Context,
	id int,
) (*domain.User, error) {
	query := `
	SELECT users.id, users.name, users.age 
	FROM myapp.users 
	WHERE users.id = $1
	`

	var u domain.User

	err := ur.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Age,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return &u, UserNotFound
	}
	if err != nil {
		return &u, DatabaseError
	}
	return &u, nil
}
