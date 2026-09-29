package userhttp

import (
	"errors"
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name *string `json:"name" validate:"required,min=2,max=255"`
	Age  *int    `json:"age" validate:"required,gte=0,lte=120"`
}

type CreateUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ToCreateUserResponse(u *domain.User) CreateUserResponse {
	return CreateUserResponse{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
	}
}

func (h *HTTPHandler) Create(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(h.Logger, w)
	var userCreateRequest CreateUserRequest

	if err := request.DecodeAndValidateRequest(r, &userCreateRequest); err != nil {
		if errors.Is(err, request.ErrValidationFailed) {
			responseHandler.ErrorResponse(errs.Invalid("validation failed", err))
			return
		} else if errors.Is(err, request.ErrExtraData) {
			responseHandler.ErrorResponse(errs.Invalid("extra data", err))
			return
		} else if errors.Is(err, request.ErrInvalidRequest) {
			responseHandler.ErrorResponse(errs.Invalid("invalid request", err))
			return
		}

		responseHandler.ErrorResponse(err)
		return
	}

	u, err := h.s.Create(r.Context(), *userCreateRequest.Name, *userCreateRequest.Age)
	if err != nil {
		responseHandler.ErrorResponse(err)
		return
	}

	createUserResponse := ToCreateUserResponse(u)

	responseHandler.JSONResponse(createUserResponse, http.StatusCreated)
}
