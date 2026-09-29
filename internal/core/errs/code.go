package errs

type Code string

const (
	CodeInvalid      Code = "invalid" // валидация / плохой ввод
	CodeNotFound     Code = "not_found"
	CodeConflict     Code = "conflict" // уникальный constraint, optimistic lock
	CodeUnauthorized Code = "unauthorized"
	CodeForbidden    Code = "forbidden"
	CodeUnavailable  Code = "unavailable" // БД недоступна, таймаут
	CodeInternal     Code = "internal"    // всё остальное непредвиденное
)
