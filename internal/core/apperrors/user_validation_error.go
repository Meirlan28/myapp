package apperrors

import "net/http"

type UserValidationError struct {
	httpCode int
	message  string
	err      error
}

func (e *UserValidationError) Error() string   { return e.message }
func (e *UserValidationError) Unwrap() error   { return e.err }
func (e *UserValidationError) HTTPStatus() int { return e.httpCode }
func (e *UserValidationError) Message() string {
	return e.message
}

func NewUserValidationError(err error, mess string) *UserValidationError {
	return &UserValidationError{httpCode: http.StatusBadRequest, message: mess, err: err}
}
