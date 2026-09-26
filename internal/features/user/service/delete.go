package service

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
)

func (us *UserService) Delete(ctx context.Context, id int) apperrors.AppError {
	err := us.Ur.DeleteByID(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return apperrors.NewDatabaseError(err, err.Error())
		case errors.Is(err, repository.UserNotFound):
			return apperrors.NewNotFoundError(err, err.Error())
		default:
			return apperrors.NewInternalServerError(err, err.Error())
		}
	}
	return nil
}
