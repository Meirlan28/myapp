package apperrors

import "net/http"

type NotFoundError struct {
	httpCode int
	message  string
	err      error
}

func (e *NotFoundError) Error() string         { return e.message }
func (e *InternalServerError) Unwrap() error   { return e.err }
func (e *InternalServerError) HTTPStatus() int { return e.httpCode }
func (e *InternalServerError) Message() string {
	return e.message
}

func NewInternalServerError(err error, mess string) *InternalServerError {
	return &InternalServerError{httpCode: http.StatusInternalServerError, message: mess, err: err}
}
