package http

import (
	"context"
	"log/slog"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/server"
	"github.com/Meirlan28/myapp/internal/features/user/service"
)

type UserHTTPHandler struct {
	userService UserService
	Logger      *slog.Logger
}

type UserService interface {
	FindAll(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, apperrors.AppError)
	CountByAge(ctx context.Context, minAge *int, maxAge *int) (int, apperrors.AppError)
	FindById(ctx context.Context, id int) (*domain.User, apperrors.AppError)
	Save(ctx context.Context, name string, age int) (*domain.User, apperrors.AppError)
	Update(ctx context.Context, id int, userUpdate domain.UserUpdate) (*domain.User, apperrors.AppError)
	Delete(ctx context.Context, id int) apperrors.AppError
}

func NewUserHTTPHandler(us *service.UserService, logger *slog.Logger) *UserHTTPHandler {
	return &UserHTTPHandler{us, logger}
}

func (h *UserHTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "POST",
			Path:    "/users",
			Handler: h.CreateUser,
		},
		{
			Method:  "DELETE",
			Path:    "/users/{id}",
			Handler: h.DeleteUser,
		},
		{
			Method:  "GET",
			Path:    "/users",
			Handler: h.GetUsers,
		},
		{
			Method:  "GET",
			Path:    "/users/{id}",
			Handler: h.FindUser,
		},
		{
			Method:  "PUT",
			Path:    "/users/{id}",
			Handler: h.UpdateUser,
		},
	}
}
