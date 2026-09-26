package http

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
	"github.com/Meirlan28/myapp/internal/core/transport/http/request"
	"github.com/Meirlan28/myapp/internal/core/transport/http/response"
)

func (uh *UserHTTPHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	responseHandler := response.NewHTTPResponseHandler(uh.Logger, w)
	var appErr apperrors.AppError
	id, err := request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to get id from path:"))
		return
	}

	appErr = uh.userService.Delete(r.Context(), id)
	if appErr != nil {
		responseHandler.ErrorResponse(apperrors.NewBadRequestError(err, "failed to delete user:"))
		return
	}
	responseHandler.NoContentResponse()
}
