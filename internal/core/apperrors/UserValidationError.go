package apperrors

import "net/http"

type UserValidationError struct {
	Code    int
	Message string
	Err     error
}

func (e *UserValidationError) Error() string { return e.Message }
func (e *UserValidationError) Unwrap() error { return e.Err }
func (e *UserValidationError) GetCode() int  { return e.Code }

func NewUserValidationError(err error) *UserValidationError {
	return &UserValidationError{Code: http.StatusBadRequest, Message: "User Validation error", Err: err}
}
