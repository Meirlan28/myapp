package apperrors

type AppError interface {
	error
	Unwrap() error
	HTTPStatus() int
	Message() string
}
