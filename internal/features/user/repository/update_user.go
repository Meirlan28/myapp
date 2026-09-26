package repository

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
)

func (ur *UserRepository) Update(
	ctx context.Context,
	id int,
	u *domain.User,
) (*domain.User, error) {
	ur.logger.Info("updating user")
	query := `
	UPDATE myapp.users 
	SET name = $1, age = $2 
	WHERE id = $3
	`
	queryResult, err := ur.db.Exec(
		ctx,
		query,
		u.Name,
		u.Age,
		id,
	)
	if err != nil {
		ur.logger.Error(err.Error())
		return &domain.User{}, DatabaseError
	}
	if queryResult.RowsAffected() == 0 {
		ur.logger.Warn("user not found")
		return &domain.User{}, UserNotFound
	}
	ur.logger.Info("user updated")
	return ur.FindByID(ctx, id)
}
