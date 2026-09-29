package userhttp

import (
	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"

	"net/http"
)

func (h *HTTPHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(h.Logger, w)
	id, err := request.IntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid id", err))
		return
	}

	u, err := h.s.GetByID(r.Context(), id)
	if err != nil {
		responseHandler.ErrorResponse(err)
		return
	}

	responseHandler.JSONResponse(u, http.StatusOK)
}
