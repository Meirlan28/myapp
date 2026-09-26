package service

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
)

func (us *UserService) Update(ctx context.Context, id int, userUpdate domain.UserUpdate) (*domain.User, apperrors.AppError) {
	user, err := us.Ur.FindByID(ctx, id)
	if err != nil {
		return &domain.User{}, apperrors.NewUserValidationError(err, err.Error())
	}

	if err := user.ApplyUpdate(userUpdate); err != nil {
		return &domain.User{}, apperrors.NewUserValidationError(err, err.Error())
	}

	updatedUser, err := us.Ur.Update(ctx, id, user)
	if err != nil {
		return &domain.User{}, apperrors.NewUserValidationError(err, err.Error())
	}

	return updatedUser, nil
}
