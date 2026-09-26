package http

import (
	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"

	"net/http"
)

func (uh *UserHTTPHandler) FindUser(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(uh.Logger, w)
	var appErr apperrors.AppError
	id, err := request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get id from path:"))
		return
	}

	u, appErr := uh.userService.FindById(r.Context(), id)
	if appErr != nil {
		responseHandler.ErrorResponse(appErr)
		return
	}

	responseHandler.JSONResponse(u, http.StatusOK)
}
