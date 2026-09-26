package service

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
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
		us.logger.Error(err.Error())
		switch {
		case errors.Is(err, repository.DatabaseError):
			return user, apperrors.NewDatabaseError(err, "database error:")
		case errors.Is(err, repository.UserNotFound):
			return user, apperrors.NewNotFoundError(err, "user not found:")
		case errors.Is(err, domain.ErrEmptyName) || errors.Is(err, domain.ErrInvalidName) || errors.Is(err, domain.ErrInvalidAge):
			return user, apperrors.NewUserValidationError(err, "invalid user data:")
		default:
			return user, apperrors.NewInternalServerError(err, "internal server error:")
		}
	}

	return updatedUser, nil
}
