package service

import (
	"context"
	"myapp/internal/core/domains/user"
)

type Repository interface {
	FindAll(ctx context.Context, minAge int, maxAge int, limit int, offset int) ([]user.User, error)
	FindByID(ctx context.Context, id int) (*user.User, error)
	Save(ctx context.Context, u *user.User) (*user.User, error)
	Update(ctx context.Context, id int, u *user.User) (*user.User, error)
	DeleteByID(ctx context.Context, id int) error
	CountByAge(ctx context.Context, minAge int, MaxAge int) (int, error)
}
