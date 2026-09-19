package service

import (
	"myapp/internal/core/domains/user"
)

type Repository interface {
	FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, error)
	FindByID(id int) (user.User, error)
	Create(name string, age int) (user.User, error)
	Update(id int, name string, age int) (user.User, error)
	DeleteByID(id int) error
	CountByAge(minAge int, MaxAge int) (int, error)
}
