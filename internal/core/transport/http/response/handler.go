package response

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/apperrors"
)

type HTTPResponseHandler struct {
	log *slog.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(
	log *slog.Logger,
	rw http.ResponseWriter,
) *HTTPResponseHandler {
	return &HTTPResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HTTPResponseHandler) JSONResponse(
	responseBody any,
	statusCode int,
) {
	h.rw.WriteHeader(statusCode)
	h.rw.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(h.rw).Encode(responseBody); err != nil {
		h.log.Error("write HTTP response", err)
	}
}

func (h *HTTPResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *HTTPResponseHandler) ErrorResponse(err apperrors.AppError) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(err.HTTPStatus())

	if err := json.NewEncoder(h.rw).Encode(ErrorResponse{err.Error(), err.Message()}); err != nil {
		h.log.Error("write HTTP error response", err)
	}
}
