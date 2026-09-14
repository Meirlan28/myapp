package repository

import (
	"encoding/json"
	"errors"
	"log/slog"
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"
	"os"

	"github.com/google/uuid"
)

const MinLimit = 1
const MaxLimit = 100

type FileUserRepository struct {
	fileName string
	users    []user.User
	logger   *slog.Logger
}

func New(fileName string, logger *slog.Logger) *FileUserRepository {
	return &FileUserRepository{fileName, []user.User{}, logger}
}

func (ur *FileUserRepository) saveToFile() apperrors.AppError {
	data, err := json.MarshalIndent(ur.users, "", "   ")
	if err != nil {
		return apperrors.NewInternalServerError(err)
	}
	err = os.WriteFile(ur.fileName, data, 0o600)
	if err != nil {
		return apperrors.NewInternalServerError(err)
	}
	return nil
}

func (ur *FileUserRepository) Load() apperrors.AppError {
	data, err := os.ReadFile(ur.fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ur.users = []user.User{}
			return nil
		}

		return apperrors.NewInternalServerError(err)
	}
	err = json.Unmarshal(data, &ur.users)
	if err != nil {
		if len(data) == 0 {
			return nil
		}
		return apperrors.NewInternalServerError(err)
	}
	return nil
}

func (ur *FileUserRepository) FindAll(minAge int, maxAge int, limit int, offset int) ([]user.User, apperrors.AppError) {
	var filteredUsers []user.User

	if minAge < user.MinAge || user.MaxAge < minAge {
		ur.logger.Info("validation error for min_age",
			"min_age", minAge)
		return nil, apperrors.NewBadRequestError(InvalidAgeError)
	}

	if maxAge < user.MinAge || user.MaxAge < maxAge {
		ur.logger.Info("validation error for max_age",
			"max_age", maxAge)
		return nil, apperrors.NewBadRequestError(InvalidAgeError)
	}

	if minAge > maxAge {
		ur.logger.Info("validation error for max_age",
			"max_age", maxAge, "min_age", minAge)
		return nil, apperrors.NewBadRequestError(InvalidAgeError)
	}

	ur.logger.Info("validation finished",
		"max_age", maxAge, "min_age", minAge)

	if limit < MinLimit || MaxLimit < limit {
		return nil, apperrors.NewBadRequestError(InvalidLimitError)
	}

	if offset < 0 {
		return nil, apperrors.NewBadRequestError(InvalidOffsetError)
	}

	limitIndex := 0
	offsetIndex := 0
	for _, u := range ur.users {
		if limitIndex >= limit {
			break
		}
		if minAge <= u.Age && u.Age <= maxAge {
			if offsetIndex < offset {
				offsetIndex++
				continue
			}
			filteredUsers = append(filteredUsers, u)
			limitIndex++
			offsetIndex++
		}
	}

	return filteredUsers, nil
}

func (ur *FileUserRepository) FindByID(id uuid.UUID) (user.User, apperrors.AppError) {
	for i := range ur.users {
		if ur.users[i].Id == id {
			return ur.users[i], nil
		}
	}
	return user.User{}, apperrors.NewNotFoundError(UserNotFound)
}

func (ur *FileUserRepository) DeleteByID(id uuid.UUID) apperrors.AppError {
	for i := range ur.users {
		if ur.users[i].Id == id {
			deletedUser := ur.users[i]
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			err := ur.saveToFile()
			if err != nil {
				ur.users = append(
					ur.users[:i],
					append([]user.User{deletedUser}, ur.users[i:]...)...,
				)
				return apperrors.NewInternalServerError(err)
			}
			return nil
		}
	}
	return apperrors.NewNotFoundError(UserNotFound)
}

func (ur *FileUserRepository) Update(id uuid.UUID, name string, age int) (user.User, apperrors.AppError) {
	var userUpdate user.User
	for i, u := range ur.users {
		if u.Id == id {
			preUpdatedUser := ur.users[i]
			ur.users[i].Name = name
			err := ur.users[i].SetAge(age)
			if err != nil {
				ur.users[i].Name = preUpdatedUser.Name
				_ = ur.users[i].SetAge(preUpdatedUser.Age)
				return user.User{}, apperrors.NewBadRequestError(InvalidAgeError)
			}
			userUpdate = ur.users[i]

			err = ur.saveToFile()
			if err != nil {
				ur.users[i].Name = preUpdatedUser.Name
				return user.User{}, apperrors.NewInternalServerError(err)
			}
			return userUpdate, nil
		}
	}
	return user.User{}, apperrors.NewNotFoundError(UserNotFound)
}

func (ur *FileUserRepository) Create(name string, age int) (user.User, apperrors.AppError) {
	u, err := user.NewUser(name, age)
	if err != nil {
		return user.User{}, apperrors.NewInternalServerError(err)
	}

	ur.users = append(ur.users, u)

	err = ur.saveToFile()
	if err != nil {
		ur.users = ur.users[:len(ur.users)-1]
		return user.User{}, apperrors.NewInternalServerError(err)
	}

	return u, nil
}

func (ur *FileUserRepository) CountByAge(minAge int, maxAge int) int {
	count := 0
	for _, u := range ur.users {
		if minAge <= u.Age && u.Age <= maxAge {
			count++
		}
	}

	return count
}
