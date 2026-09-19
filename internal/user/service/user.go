package service

import (
	"errors"
	"fmt"
	"log/slog"
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"
	"myapp/internal/user/repository"
)

type UserService struct {
	Ur     Repository
	logger *slog.Logger
}

func New(ur Repository, logger *slog.Logger) *UserService {
	return &UserService{ur, logger}
}

func (us *UserService) Create(name string, age int) (user.User, apperrors.AppError) {
	u, err := us.Ur.Create(name, age)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return u, apperrors.NewDatabaseError(err)
		default:
			return u, apperrors.NewInternalServerError(err)
		}
	}
	return u, nil
}

func (us *UserService) Delete(id int) apperrors.AppError {
	err := us.Ur.DeleteByID(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return apperrors.NewDatabaseError(err)
		case errors.Is(err, repository.UserNotFound):
			return apperrors.NewNotFoundError(err)
		default:
			return apperrors.NewInternalServerError(err)
		}
	}
	return nil
}

func (us *UserService) FindById(id int) (user.User, apperrors.AppError) {
	us.logger.Info("finding user by id")
	u, err := us.Ur.FindByID(id)
	if err != nil {
		us.logger.Error(err.Error())
		switch {
		case errors.Is(err, repository.DatabaseError):
			return u, apperrors.NewDatabaseError(err)
		case errors.Is(err, repository.UserNotFound):
			return u, apperrors.NewNotFoundError(err)
		default:
			return u, apperrors.NewInternalServerError(err)
		}
	}
	return u, nil
}

func (us *UserService) FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, apperrors.AppError) {
	if minAge < user.MinAge || user.MaxAge < minAge {
		us.logger.Info("validation error for min_age",
			"min_age", minAge)
		return nil, apperrors.NewBadRequestError(fmt.Errorf("invalid age: %d", minAge))
	}

	if maxAge < user.MinAge || user.MaxAge < maxAge {
		us.logger.Info("validation error for max_age",
			"max_age", maxAge)
		return nil, apperrors.NewBadRequestError(fmt.Errorf("invalid age: %d", minAge))
	}

	if minAge > maxAge {
		us.logger.Info("validation error for max_age",
			"max_age", maxAge, "min_age", minAge)
		return nil, apperrors.NewBadRequestError(fmt.Errorf("invalid age range: %d-%d", minAge, maxAge))
	}

	us.logger.Info("validation finished",
		"max_age", maxAge, "min_age", minAge)

	if limit < user.MinLimit || user.MaxLimit < limit {
		return nil, apperrors.NewBadRequestError(fmt.Errorf("invalid limit: %d", limit))
	}

	if offset < 0 {
		return nil, apperrors.NewBadRequestError(fmt.Errorf("invalid offset: %d", offset))
	}

	users, err := us.Ur.FindAll(minAge, maxAge, limit, offset)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return users, apperrors.NewDatabaseError(err)
		default:
			return users, apperrors.NewInternalServerError(err)
		}
	}
	return users, nil
}

func (us *UserService) Update(id int, name string, age int) (user.User, apperrors.AppError) {
	u, err := us.Ur.Update(id, name, age)
	if err != nil {
		us.logger.Error(err.Error())
		switch {
		case errors.Is(err, repository.DatabaseError):
			return u, apperrors.NewDatabaseError(err)
		case errors.Is(err, repository.UserNotFound):
			return u, apperrors.NewNotFoundError(err)
		default:
			return u, apperrors.NewInternalServerError(err)
		}
	}
	return u, nil
}

func (us *UserService) CountByAge(minAge int, MaxAge int) (int, apperrors.AppError) {
	count, err := us.Ur.CountByAge(minAge, MaxAge)
	if err != nil {
		switch {
		case errors.Is(err, repository.DatabaseError):
			return 0, apperrors.NewDatabaseError(err)
		case errors.Is(err, repository.UserNotFound):
			return 0, apperrors.NewNotFoundError(err)
		default:
			return 0, apperrors.NewInternalServerError(err)
		}
	}
	return count, nil
}
