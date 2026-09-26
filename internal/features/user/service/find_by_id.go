package service

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
)

func (us *UserService) FindById(ctx context.Context, id int) (*domain.User, apperrors.AppError) {
	us.logger.Info("finding user by id")
	u, err := us.Ur.FindByID(ctx, id)
	if err != nil {
		us.logger.Error(err.Error())
		switch {
		case errors.Is(err, repository.DatabaseError):
			return u, apperrors.NewDatabaseError(err, err.Error())
		case errors.Is(err, repository.UserNotFound):
			return u, apperrors.NewInternalServerError(err, err.Error())
		default:
			return u, apperrors.NewInternalServerError(err, err.Error())
		}
	}
	return u, nil
}
