package apperrors

import "net/http"

type InternalServerError struct {
	Code    int
	Message string
	Err     error
}

func (e *InternalServerError) Error() string { return e.Message }
func (e *InternalServerError) Unwrap() error { return e.Err }
func (e *InternalServerError) GetCode() int  { return e.Code }

func NewInternalServerError(err error) *InternalServerError {
	return &InternalServerError{Code: http.StatusInternalServerError, Message: "Internal Server Error", Err: err}
}
