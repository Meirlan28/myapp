package apperrors

type AppError interface {
	Error() string
	Unwrap() error
	GetCode() int
}
