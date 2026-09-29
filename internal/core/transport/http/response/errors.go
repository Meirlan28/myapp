package response

import (
	"net/http"

	"github.com/Meirlan28/myapp/internal/core/errs"
)

type ErrorResponse struct {
	Code    errs.Code `json:"code"`
	Message string    `json:"message"`
}

// Code to HTTP Status Code
func httpStatus(c errs.Code) int {
	switch c {
	case errs.CodeInvalid:
		return http.StatusBadRequest
	case errs.CodeNotFound:
		return http.StatusNotFound
	case errs.CodeConflict:
		return http.StatusConflict
	case errs.CodeUnauthorized:
		return http.StatusUnauthorized
	case errs.CodeForbidden:
		return http.StatusForbidden
	case errs.CodeUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
