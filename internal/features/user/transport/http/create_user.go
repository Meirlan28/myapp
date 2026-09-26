package http

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	Name *string `json:"name" validate:"required",min=2,max=255`
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

func (uh *UserHTTPHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(uh.Logger, w)
	var appErr apperrors.AppError
	var userCreateRequest CreateUserRequest

	err := request.DecodeAndValidateRequest(r, &userCreateRequest)
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to decode request body:"))
		return
	}

	u, appErr := uh.userService.Save(r.Context(), *userCreateRequest.Name, *userCreateRequest.Age)
	if appErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to create user:"))
		return
	}

	createUserResponse := ToCreateUserResponse(u)

	responseHandler.JSONResponse(createUserResponse, http.StatusCreated)
}
