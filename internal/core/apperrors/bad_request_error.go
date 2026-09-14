package apperrors

import "net/http"

type BadRequestError struct {
	Code    int
	Message string
	Err     error
}

func (e *BadRequestError) Error() string { return e.Message }
func (e *BadRequestError) Unwrap() error { return e.Err }
func (e *BadRequestError) GetCode() int  { return e.Code }

func NewBadRequestError(err error) *BadRequestError {
	return &BadRequestError{Code: http.StatusBadRequest, Message: "Bad Request", Err: err}
}
