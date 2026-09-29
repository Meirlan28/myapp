package service

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
)

func (us *UserService) List(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, error) {

	if limit != nil && *limit < 0 {
		return nil, errs.Invalid("invalid limit", ErrInvalidLimit)
	}

	if offset != nil && *offset < 0 {
		return nil, errs.Invalid("invalid offset", ErrInvalidOffset)
	}

	users, err := us.r.List(ctx, minAge, maxAge, limit, offset)
	if err != nil {
		return nil, err
	}
	return users, nil
}
