package service

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/postgres"
)

func (us *UserService) Update(ctx context.Context, id int, userUpdate domain.UserUpdate) (*domain.User, error) {
	user, err := us.r.GetByID(ctx, id)
	if err != nil {
		return nil, postgres.MapError(err)
	}

	if err := user.ApplyUpdate(userUpdate); err != nil {
		return nil, errs.Invalid("validation error", err)
	}

	updatedUser, err := us.r.Update(ctx, id, user)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}
