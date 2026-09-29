package errs

type Error struct {
	Code    Code
	Message string // Безопасно показывать клиенту
	Err     error  // Настоящая причина — только для логов, никогда не в JSON-ответ
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Err }

func Invalid(msg string, cause error) *Error      { return &Error{CodeInvalid, msg, cause} }
func NotFound(msg string, cause error) *Error     { return &Error{CodeNotFound, msg, cause} }
func Conflict(msg string, cause error) *Error     { return &Error{CodeConflict, msg, cause} }
func Unauthorized(msg string, cause error) *Error { return &Error{CodeUnauthorized, msg, cause} }
func Forbidden(msg string, cause error) *Error    { return &Error{CodeForbidden, msg, cause} }
func Unavailable(msg string, cause error) *Error  { return &Error{CodeUnavailable, msg, cause} }
func Internal(msg string, cause error) *Error     { return &Error{CodeInternal, msg, cause} }
