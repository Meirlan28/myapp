package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/errs"
)

type ResponseHandler struct {
	log *slog.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *slog.Logger,
	rw http.ResponseWriter,
) *ResponseHandler {
	return &ResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *ResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)
	h.rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error(
			"failed to encode response body",
			"error", err,
		)
	}
}

func (h *ResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *ResponseHandler) ErrorResponse(err error) {
	var appErr *errs.Error
	if !errors.As(err, &appErr) {
		appErr = errs.Internal("internal error", err)
	}

	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(httpStatus(appErr.Code))

	if err := json.NewEncoder(h.rw).Encode(ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	}); err != nil {
		h.log.Error("write HTTP error response", "error", err)
	}
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	err := errs.Internal("internal error", fmt.Errorf("%v", p))
	h.log.Error(msg, "error", err)

	h.ErrorResponse(
		errs.Internal("internal error", err),
	)
}
