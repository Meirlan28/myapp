package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"myapp/internal/user"

	"github.com/google/uuid"
)

const MIN_LIMIT = 1
const MAX_LIMIT = 100

type UserRepository struct {
	fileName string
	users    []user.User
}

func New(fileName string) *UserRepository {
	return &UserRepository{fileName: fileName, users: []user.User{}}
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
	if min_age < user.MinAge || user.MaxAge < max_age {
		return nil, ValidationError
	}

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
		if offset_index < offset {
			offset_index++
			continue
		}
		if min_age <= u.Age && u.Age <= max_age {
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
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			err := ur.saveToFile()
			if err != nil {
				return FileError
			}
			return nil
		}
	}
	return UserNotFound
}

func (ur *UserRepository) UpdateAge(id uuid.UUID, age int) (user.User, error) {
	var userUpdate user.User
	for i, u := range ur.users {
		if u.Id == id {
			err := ur.users[i].SetAge(age)
			if err != nil {
				return u, err
			}
			userUpdate = ur.users[i]

			err = ur.saveToFile()
			if err != nil {
				return user.User{}, FileError
			}
			return userUpdate, nil
		}
	}
	return user.User{}, UserNotFound
}

func (ur *UserRepository) UpdateName(id uuid.UUID, name string) (user.User, error) {
	var userUpdate user.User
	for i, u := range ur.users {
		if u.Id == id {
			ur.users[i].Name = name
			userUpdate = ur.users[i]

			err := ur.saveToFile()
			if err != nil {
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

var ClientError = errors.New("bad request")

var UserNotFound = fmt.Errorf("%w: user not found", ClientError)

var InvalidAge = fmt.Errorf("%w: invalid age", ClientError)

var JsonParsingError = fmt.Errorf("%w: failed to parse to json", ClientError)

var FileError = fmt.Errorf("%w: file error", ServerError)

var ValidationError = fmt.Errorf("%w: validation error", ClientError)
