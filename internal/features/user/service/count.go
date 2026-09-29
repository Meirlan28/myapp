package service

import (
	"context"
)

func (us *UserService) Count(ctx context.Context, minAge *int, maxAge *int) (int, error) {
	count, err := us.r.Count(ctx, minAge, maxAge)
	return count, err
}
