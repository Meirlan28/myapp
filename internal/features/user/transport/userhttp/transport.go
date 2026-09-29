package userhttp

import (
	"context"
	"log/slog"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/server"
	"github.com/Meirlan28/myapp/internal/features/user/service"
)

type HTTPHandler struct {
	s      UserService
	Logger *slog.Logger
}

type UserService interface {
	List(ctx context.Context, minAge *int, maxAge *int, limit *int, offset *int) ([]domain.User, error)
	Count(ctx context.Context, minAge *int, maxAge *int) (int, error)
	GetByID(ctx context.Context, id int) (*domain.User, error)
	Create(ctx context.Context, name string, age int) (*domain.User, error)
	Update(ctx context.Context, id int, userUpdate domain.UserUpdate) (*domain.User, error)
	Delete(ctx context.Context, id int) error
}

func NewHTTPHandler(us *service.UserService, logger *slog.Logger) *HTTPHandler {
	return &HTTPHandler{us, logger}
}

func (h *HTTPHandler) Routes() []server.Route {
	return []server.Route{
		{
			Method:  "POST",
			Path:    "/users",
			Handler: h.Create,
		},
		{
			Method:  "DELETE",
			Path:    "/users/{id}",
			Handler: h.Delete,
		},
		{
			Method:  "GET",
			Path:    "/users",
			Handler: h.List,
		},
		{
			Method:  "GET",
			Path:    "/users/{id}",
			Handler: h.GetByID,
		},
		{
			Method:  "PUT",
			Path:    "/users/{id}",
			Handler: h.Update,
		},
	}
}
