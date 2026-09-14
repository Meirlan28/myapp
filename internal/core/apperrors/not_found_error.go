package apperrors

import "net/http"

type NotFoundError struct {
	Code    int
	Message string
	Err     error
}

func (e *NotFoundError) Error() string { return e.Message }
func (e *NotFoundError) Unwrap() error { return e.Err }
func (e *NotFoundError) GetCode() int  { return e.Code }

func NewNotFoundError(err error) *InternalServerError {
	return &InternalServerError{Code: http.StatusNotFound, Message: "Not Found Error", Err: err}
}
