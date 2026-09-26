package apperrors

import "net/http"

type NotFoundError struct {
	httpCode int
	message  string
	err      error
}

func (e *NotFoundError) Error() string   { return e.message }
func (e *NotFoundError) Unwrap() error   { return e.err }
func (e *NotFoundError) HTTPStatus() int { return e.httpCode }
func (e *NotFoundError) Message() string {
	return e.message
}

func NewNotFoundError(err error, mess string) *NotFoundError {
	return &NotFoundError{httpCode: http.StatusInternalServerError, message: mess, err: err}
}
