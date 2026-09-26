package http

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	domain "github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
	"github.com/Meirlan28/myapp/internal/core/transport/http/types"
)

type UpdateUserRequest struct {
	Name types.Nullable[string] `json:"name"`
	Age  types.Nullable[int]    `json:"age"`
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
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get id from path:"))
		return
	}

	var userUpdateRequest UpdateUserRequest
	err = request.DecodeAndValidateRequest(r, &userUpdateRequest)
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to decode request body:"))
		return
	}

	userUpdate := userPatchFromRequest(userUpdateRequest)

	u, AppErr := uh.userService.FindById(r.Context(), id)
	if AppErr != nil {
		responseHandler.ErrorResponse(AppErr)
		return
	}

	u, appErr = uh.userService.Update(r.Context(), id, userUpdate)
	if appErr != nil {
		responseHandler.ErrorResponse(appErr)
		return
	}

	responseHandler.JSONResponse(ToUpdateUserResponse(u), http.StatusOK)
}

func userPatchFromRequest(request UpdateUserRequest) domain.UserUpdate {
	return domain.NewUserUpdate(
		request.Name.ToDomain(),
		request.Age.ToDomain(),
	)
}
