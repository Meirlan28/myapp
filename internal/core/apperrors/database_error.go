package apperrors

import "net/http"

type DatabaseError struct {
	Code    int
	Message string
	Err     error
}

func (e *DatabaseError) Error() string { return e.Message }
func (e *DatabaseError) Unwrap() error { return e.Err }
func (e *DatabaseError) GetCode() int  { return e.Code }

func NewDatabaseError(err error) *DatabaseError {
	return &DatabaseError{Code: http.StatusInternalServerError, Message: "Database error", Err: err}
}
