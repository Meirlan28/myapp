package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
)

func (ur *UserRepository) Save(
	ctx context.Context,
	u *domain.User,
) (*domain.User, error) {
	query := `
	INSERT INTO myapp.users (name, age) 
	VALUES ($1, $2) 
	RETURNING id, name, age
	`

	err := ur.db.QueryRow(
		ctx,
		query,
		u.Name,
		u.Age,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Age,
	)
	if err != nil {
		return u, DatabaseError
	}
	return u, nil
}
