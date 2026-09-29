package userhttp

import (
	"errors"
	"net/http"

	domain "github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
	"github.com/Meirlan28/myapp/internal/core/transport/http/types"
)

type UpdateUserRequest struct {
	Name types.Optional[string] `json:"name"`
	Age  *int                   `json:"age"`
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

func (h *HTTPHandler) Update(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(h.Logger, w)

	id, err := request.IntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid id", err))
		return
	}

	var userUpdateRequest UpdateUserRequest
	if err := request.DecodeAndValidateRequest(r, &userUpdateRequest); err != nil {
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

	userUpdate := userPatchFromRequest(userUpdateRequest)

	u, err := h.s.Update(r.Context(), id, userUpdate)
	if err != nil {
		responseHandler.ErrorResponse(err)
		return
	}

	responseHandler.JSONResponse(ToUpdateUserResponse(u), http.StatusOK)
}

func userPatchFromRequest(request UpdateUserRequest) domain.UserUpdate {
	return domain.NewUserUpdate(
		request.Name.ToDomain(),
		request.Age,
	)
}
