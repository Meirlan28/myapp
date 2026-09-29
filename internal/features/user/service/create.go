package service

import (
	"context"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
)

func (us *UserService) Create(ctx context.Context, name string, age int) (*domain.User, error) {
	u, err := domain.NewUser(name, age)
	if err != nil {
		return nil, errs.Invalid("Invalid data", err)
	}

	return us.r.Create(ctx, u)
}
