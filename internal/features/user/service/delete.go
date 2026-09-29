package service

import (
	"context"
)

func (us *UserService) Delete(ctx context.Context, id int) error {
	return us.r.Delete(ctx, id)
}
