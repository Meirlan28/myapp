package userhttp

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/errs"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

func (h *HTTPHandler) Delete(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(h.Logger, w)
	id, err := request.IntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(errs.Invalid("invalid id", err))
		return
	}

	if err := h.s.Delete(r.Context(), id); err != nil {
		responseHandler.ErrorResponse(err)
		return
	}

	responseHandler.NoContentResponse()
}
