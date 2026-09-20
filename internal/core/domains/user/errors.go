package user

import "errors"

var (
	ErrEmptyName   = errors.New("empty name")
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidAge  = errors.New("invalid age")
)
