package service

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
)

func (us *UserService) Update(ctx context.Context, id int, name string, age int) (*domain.User, apperrors.AppError) {
	u, err := domain.New(name, age)
	if err != nil {
		return &domain.User{}, apperrors.NewUserValidationError(err, err.Error())
	}
	u, err = us.Ur.Update(ctx, id, u)
	if err != nil {
		us.logger.Error(err.Error())
		switch {
		case errors.Is(err, repository.DatabaseError):
			return u, apperrors.NewDatabaseError(err, err.Error())
		case errors.Is(err, repository.UserNotFound):
			return u, apperrors.NewNotFoundError(err, err.Error())
		case errors.Is(err, domain.ErrEmptyName) || errors.Is(err, domain.ErrInvalidName) || errors.Is(err, domain.ErrInvalidAge):
			return u, apperrors.NewUserValidationError(err, err.Error())
		default:
			return u, apperrors.NewInternalServerError(err, err.Error())
		}
	}
	return u, nil
}
