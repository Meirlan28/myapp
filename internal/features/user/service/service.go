package service

import (
	"context"
	"log/slog"

	"github.com/Meirlan28/myapp/internal/core/domain"
)

type Repository interface {
	FindAll(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, error)
	FindByID(ctx context.Context, id int) (*domain.User, error)
	Save(ctx context.Context, u *domain.User) (*domain.User, error)
	Update(ctx context.Context, id int, u *domain.User) (*domain.User, error)
	DeleteByID(ctx context.Context, id int) error
	CountByAge(ctx context.Context, minAge *int, MaxAge *int) (int, error)
}
type UserService struct {
	Ur     Repository
	logger *slog.Logger
}

func New(ur Repository, logger *slog.Logger) *UserService {
	return &UserService{ur, logger}
}
