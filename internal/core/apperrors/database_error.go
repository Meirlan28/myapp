package apperrors

import "net/http"

type DatabaseError struct {
	httpCode int
	message  string
	err      error
}

func (e *DatabaseError) Error() string   { return e.message }
func (e *DatabaseError) Unwrap() error   { return e.err }
func (e *DatabaseError) HTTPStatus() int { return e.httpCode }
func (e *DatabaseError) Message() string {
	return e.message
}

func NewDatabaseError(err error, mess string) *DatabaseError {
	return &DatabaseError{httpCode: http.StatusInternalServerError, message: mess, err: err}
}
