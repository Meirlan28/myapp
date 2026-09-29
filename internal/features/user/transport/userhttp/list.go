package userhttp

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/domain"
	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

func (h *HTTPHandler) List(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(h.Logger, w)

	minAge, err := request.IntQueryParam(r, "min_age")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid min_age", err))
		return
	}

	maxAge, err := request.IntQueryParam(r, "max_age")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid max_age", err))
		return
	}

	limit, err := request.IntQueryParam(r, "limit")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid limit", err))
		return
	}

	offset, err := request.IntQueryParam(r, "offset")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid offset", err))
		return
	}

	users, err := h.s.List(r.Context(), minAge, maxAge, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err)
		return
	}

	count, err := h.s.Count(r.Context(), minAge, maxAge)
	if err != nil {
		responseHandler.ErrorResponse(err)
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
