package service

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
)

func (us *UserService) GetByID(ctx context.Context, id int) (*domain.User, error) {
	return us.r.GetByID(ctx, id)
}
