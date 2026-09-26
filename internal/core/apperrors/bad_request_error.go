package apperrors

import "net/http"

type BadRequestError struct {
	httpCode int
	message  string
	err      error
}

func (e *BadRequestError) Error() string   { return e.message }
func (e *BadRequestError) Unwrap() error   { return e.err }
func (e *BadRequestError) HTTPStatus() int { return e.httpCode }
func (e *BadRequestError) Message() string {
	return e.message
}

func NewBadRequestError(err error, mess string) *BadRequestError {
	return &BadRequestError{httpCode: http.StatusBadRequest, message: mess, err: err}
}
