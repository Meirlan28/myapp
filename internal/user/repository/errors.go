package repository

import "errors"

var InvalidAgeError = errors.New("invalid age")

var InvalidLimitError = errors.New("invalid limit")

var InvalidOffsetError = errors.New("invalid offset")

var UserNotFound = errors.New("user not found")
