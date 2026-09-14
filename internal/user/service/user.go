package service

import (
	"log/slog"
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"

	"github.com/google/uuid"
)

type UserService struct {
	Ur     Repository
	logger *slog.Logger
}

func New(ur Repository, logger *slog.Logger) *UserService {
	return &UserService{ur, logger}
}

func (us *UserService) Create(name string, age int) (user.User, apperrors.AppError) {
	return us.Ur.Create(name, age)
}

func (us *UserService) Delete(id uuid.UUID) apperrors.AppError {
	return us.Ur.DeleteByID(id)
}

func (us *UserService) FindById(id uuid.UUID) (user.User, apperrors.AppError) {
	return us.Ur.FindByID(id)
}

func (us *UserService) FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, apperrors.AppError) {
	return us.Ur.FindAll(minAge, maxAge, limit, offset)
}

func (us *UserService) Update(id uuid.UUID, name string, age int) (user.User, apperrors.AppError) {
	return us.Ur.Update(id, name, age)
}

func (us *UserService) CountByAge(minAge int, MaxAge int) (int, apperrors.AppError) {
	return us.Ur.CountByAge(minAge, MaxAge)
}
