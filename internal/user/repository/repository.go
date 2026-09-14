package repository

import (
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"

	"github.com/google/uuid"
)

type Repository interface {
	FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, apperrors.AppError)
	FindByID(id uuid.UUID) (user.User, apperrors.AppError)
	Create(name string, age int) (user.User, apperrors.AppError)
	Update(id uuid.UUID, name string, age int) (user.User, apperrors.AppError)
	DeleteByID(id uuid.UUID) apperrors.AppError
	CountByAge(minAge int, MaxAge int) int
}
