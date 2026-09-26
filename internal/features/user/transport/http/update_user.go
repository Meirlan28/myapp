package http

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	domain "github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

type UpdateUserRequest struct {
	Name *string `json:"name" min=2,max=255`
	Age  *int    `json:"age" gte=0,lte=120`
}

type UpdateUserResponse struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func ToUpdateUserResponse(u *domain.User) UpdateUserResponse {
	return UpdateUserResponse{
		ID:   u.ID,
		Name: u.Name,
		Age:  u.Age,
	}
}

func (uh *UserHTTPHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(uh.Logger, w)
	var appErr apperrors.AppError

	id, err := request.GetIntPathValue(r, "id")

	var userUpdateRequest UpdateUserRequest
	err = request.DecodeAndValidateRequest(r, &userUpdateRequest)
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to decode request body:"))
		return
	}

	u, AppErr := uh.userService.FindById(r.Context(), id)
	if AppErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to find user:"))
		return
	}

	u, appErr = uh.userService.Update(r.Context(), id, *userUpdateRequest.Name, *userUpdateRequest.Age)
	if appErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to update user:"))
		return
	}

	responseHandler.JSONResponse(ToUpdateUserResponse(u), http.StatusOK)
}
