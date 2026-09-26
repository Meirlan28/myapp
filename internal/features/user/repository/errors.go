package repository

import "errors"

var UserNotFound = errors.New("user not found")

var DatabaseError = errors.New("database error")
