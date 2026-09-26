package http

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

func (uh *UserHTTPHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(uh.Logger, w)
	var appErr apperrors.AppError
	minAge, err := request.GetIntQueryParam(r, "min_age")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get min_age from query params:"))
		return
	}

	maxAge, err := request.GetIntQueryParam(r, "max_age")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get max_age from query params:"))
		return
	}

	limit, err := request.GetIntQueryParam(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get limit from query params:"))
		return
	}

	offset, err := request.GetIntQueryParam(r, "offset")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get offset from query params:"))
		return
	}

	users, appErr := uh.userService.FindAll(r.Context(), minAge, maxAge, limit, offset)
	if appErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get users:"))
		return
	}
	count, appErr := uh.userService.CountByAge(r.Context(), minAge, maxAge)
	if appErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get users:"))
		return
	}

	var userPage UserPage

	if users == nil {
		userPage = UserPage{[]domain.User{}, count, limit, offset}
	} else {
		userPage = UserPage{users, count, limit, offset}
	}

	responseHandler.JSONResponse(userPage, http.StatusOK)
}
