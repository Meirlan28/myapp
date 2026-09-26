package service

import (
	"context"
	"errors"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
)

func (us *UserService) CountByAge(ctx context.Context, minAge *int, MaxAge *int) (int, apperrors.AppError) {
	count, err := us.Ur.CountByAge(ctx, minAge, MaxAge)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return 0, apperrors.NewDatabaseError(err, err.Error())
		case errors.Is(err, repository.UserNotFound):
			return 0, apperrors.NewInternalServerError(err, err.Error())
		default:
			return 0, apperrors.NewInternalServerError(err, err.Error())
		}
	}
	return count, nil
}
