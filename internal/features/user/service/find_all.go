package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
)

func (us *UserService) FindAll(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, apperrors.AppError) {

	if limit != nil && *limit < 0 {
		return nil, apperrors.AppError(apperrors.NewBadRequestError(fmt.Errorf("limit must be greater than 0"), "limit must be greater than 0"))
	}

	if offset != nil && *offset < 0 {
		return nil, apperrors.AppError(apperrors.NewBadRequestError(fmt.Errorf("offset must be greater than 0"), "offset must be greater than 0"))
	}

	us.logger.Info("validation finished",
		"max_age", maxAge, "min_age", minAge)

	users, err := us.Ur.FindAll(ctx, minAge, maxAge, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return users, apperrors.NewDatabaseError(err, err.Error())
		default:
			return users, apperrors.NewInternalServerError(err, err.Error())
		}
	}
	return users, nil
}
