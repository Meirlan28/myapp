package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"myapp/internal/user"

	"github.com/google/uuid"
)

const MIN_LIMIT = 1
const MAX_LIMIT = 100

type UserRepository struct {
	fileName string
	users    []user.User
	logger   *slog.Logger
}

func New(fileName string, logger *slog.Logger) *UserRepository {
	return &UserRepository{fileName, []user.User{}, logger}
}

func (ur *UserRepository) saveToFile() error {
	data, err := json.MarshalIndent(ur.users, "", "   ")
	if err != nil {
		return JsonParsingError
	}
	err = os.WriteFile(ur.fileName, data, 0o600)
	if err != nil {
		return FileError
	}
	return nil
}

func (ur *UserRepository) Load() error {
	data, err := os.ReadFile(ur.fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ur.users = []user.User{}
			return nil
		}

		return FileError
	}
	err = json.Unmarshal(data, &ur.users)
	if err != nil {
		if len(data) == 0 {
			return nil
		}
		return JsonParsingError
	}
	return nil
}

func (ur *UserRepository) GetAllFiltered(min_age int, max_age int, limit int, offset int) ([]user.User, error) {
	var filteredUsers []user.User

	if min_age < user.MinAge || user.MaxAge < min_age {
		ur.logger.Info("validation error for min_age",
			"min_age", min_age)
		return nil, ValidationError
	}

	if max_age < user.MinAge || user.MaxAge < max_age {
		ur.logger.Info("validation error for max_age",
			"max_age", max_age)
		return nil, ValidationError
	}

	if min_age > max_age {
		ur.logger.Info("validation error for max_age",
			"max_age", max_age, "min_age", min_age)
		return nil, ValidationError
	}

	ur.logger.Info("validation finished",
		"max_age", max_age, "min_age", min_age)

	if limit < MIN_LIMIT || MAX_LIMIT < limit {
		return nil, ValidationError
	}

	if offset < 0 {
		return nil, ValidationError
	}

	limit_index := 0
	offset_index := 0
	for _, u := range ur.users {
		if limit_index >= limit {
			break
		}
		if min_age <= u.Age && u.Age <= max_age {
			if offset_index < offset {
				offset_index++
				continue
			}
			filteredUsers = append(filteredUsers, u)
			limit_index++
			offset_index++
		}
	}

	return filteredUsers, nil
}

func (ur *UserRepository) GetAll() []user.User {
	return ur.users
}

func (ur *UserRepository) FindByID(id uuid.UUID) (user.User, error) {
	for i := range ur.users {
		if ur.users[i].Id == id {
			return ur.users[i], nil
		}
	}
	return user.User{}, UserNotFound
}

func (ur *UserRepository) DeleteByID(id uuid.UUID) error {
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
				return FileError
			}
			return nil
		}
	}
	return UserNotFound
}

func (ur *UserRepository) UpdateAge(id uuid.UUID, age int) (user.User, error) {
	var userUpdated user.User
	for i, u := range ur.users {
		if u.Id == id {
			preUpdatedUser := ur.users[i]
			err := ur.users[i].SetAge(age)
			if err != nil {
				return u, err
			}
			userUpdated = ur.users[i]

			err = ur.saveToFile()
			if err != nil {
				ur.users[i].SetAge(preUpdatedUser.Age)
				return user.User{}, FileError
			}
			return userUpdated, nil
		}
	}
	return user.User{}, UserNotFound
}

func (ur *UserRepository) UpdateName(id uuid.UUID, name string) (user.User, error) {
	var userUpdate user.User
	for i, u := range ur.users {
		if u.Id == id {
			preUpdatedUser := ur.users[i]
			ur.users[i].Name = name
			userUpdate = ur.users[i]

			err := ur.saveToFile()
			if err != nil {
				ur.users[i].Name = preUpdatedUser.Name
				return user.User{}, FileError
			}
			return userUpdate, nil
		}
	}
	return user.User{}, UserNotFound
}

func (ur *UserRepository) Create(name string, age int) (user.User, error) {
	u, err := user.NewUser(name, age)
	if err != nil {
		return user.User{}, FileError
	}

	ur.users = append(ur.users, u)

	err = ur.saveToFile()
	if err != nil {
		ur.users = ur.users[:len(ur.users)-1]
		return user.User{}, FileError
	}

	return u, nil
}

func (ur *UserRepository) CountByAge(min_age int, max_age int) int {
	count := 0
	for _, u := range ur.users {
		if min_age <= u.Age && u.Age <= max_age {
			count++
		}
	}

	return count
}

var ServerError = errors.New("internal server error")

var FileError = fmt.Errorf("%w: file error", ServerError)

var ClientError = errors.New("bad request")

var UserNotFound = fmt.Errorf("%w: user not found", ClientError)

var InvalidAge = fmt.Errorf("%w: invalid age", ClientError)

var JsonParsingError = fmt.Errorf("%w: failed to parse to json", ServerError)

var InvalidQueryParameter = fmt.Errorf("%w: invalid query parameter", ClientError)

var ValidationError = fmt.Errorf("%w: validation error", ClientError)

var InvalidUUID = fmt.Errorf("%w: invalid uuid", ClientError)

var InvalidBodyInRequest = fmt.Errorf("%w: invalid body in request", ClientError)
