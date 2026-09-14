package service

import (
	"log/slog"
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"
	"myapp/internal/user/repository"

	"github.com/google/uuid"
)

type UserService struct {
	Ur     *repository.UserRepository
	logger *slog.Logger
}

func New(ur *repository.UserRepository, logger *slog.Logger) *UserService {
	return &UserService{ur, logger}
}

func (us *UserService) CreateUser(name string, age int) (user.User, apperrors.AppError) {
	return us.Ur.Create(name, age)
}

func (us *UserService) UpdateUser(id uuid.UUID, age int) (user.User, apperrors.AppError) {
	return us.Ur.UpdateAge(id, age)
}

func (us *UserService) DeleteUser(id uuid.UUID) apperrors.AppError {
	return us.Ur.DeleteByID(id)
}

func (us *UserService) GetUser(id uuid.UUID) (user.User, apperrors.AppError) {
	return us.Ur.FindByID(id)
}

func (us *UserService) GetUsers(minAge int, maxAge int, limit int, offset int) ([]user.User, apperrors.AppError) {
	return us.Ur.FindAll(minAge, maxAge, limit, offset)
}

func (us *UserService) UpdateAge(id uuid.UUID, age int) (user.User, apperrors.AppError) {
	return us.Ur.UpdateAge(id, age)
}

func (us *UserService) UpdateName(id uuid.UUID, name string) (user.User, apperrors.AppError) {
	return us.Ur.UpdateName(id, name)
}

func (us *UserService) CountByAge(minAge int, MaxAge int) int {
	return us.Ur.CountByAge(minAge, MaxAge)
}
