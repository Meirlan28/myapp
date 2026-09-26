package http

import (
	"github.com/Meirlan28/myapp/internal/core/domain"
)

type UserPage struct {
	Users  []domain.User `json:"users"`
	Total  int           `json:"total"`
	Limit  *int          `json:"limit"`
	Offset *int          `json:"offset"`
}
