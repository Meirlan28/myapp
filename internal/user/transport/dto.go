package transport

import "myapp/internal/core/domains/user"

type UserCreateRequest struct {
	Name *string
	Age  *int
}

type UserUpdateRequest struct {
	Name *string
	Age  *int
}

type UserPage struct {
	Users  []user.User `json:"users"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
