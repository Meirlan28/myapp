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
		return err
	}
	err = os.WriteFile(ur.fileName, data, 0o600)
	if err != nil {
		return err
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

		return fmt.Errorf("read %s: %w", ur.fileName, err)
	}
	err = json.Unmarshal(data, &ur.users)
	if err != nil {
		if len(data) == 0 {
			return nil
		}
		return fmt.Errorf("invalid JSON in %s: %w", ur.fileName, err)
	}
	return nil
}

func (ur *UserRepository) GetAllFiltered(min_age int, max_age int, limit int, offset int) ([]user.User, error) {
	var filteredUsers []user.User
	if min_age < user.MinAge || user.MaxAge < max_age {
		return nil, fmt.Errorf("invalid min_age or max_age")
	}

	if limit < MIN_LIMIT || MAX_LIMIT < limit {
		return nil, fmt.Errorf("invalid limit")
	}

	if offset < 0 {
		return nil, fmt.Errorf("offset must be positive")
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
	return user.User{}, errors.New("User not found")
}

func (ur *UserRepository) DeleteByID(id uuid.UUID) error {
	for i := range ur.users {
		if ur.users[i].Id == id {
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			break
		}
	}
	err := ur.saveToFile()
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Update(id uuid.UUID, name string, age int) (user.User, error) {
	var user user.User
	for i, u := range ur.users {
		if u.Id == id {
			err := ur.users[i].SetAge(age)
			if err != nil {
				return u, err
			}
			ur.users[i].Name = name
			user = ur.users[i]
		}
	}

	err := ur.saveToFile()
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ur *UserRepository) Create(name string, age int) (user.User, error) {
	u, err := user.NewUser(name, age)
	if err != nil {
		return user.User{}, err
	}

	ur.users = append(ur.users, u)

	err = ur.saveToFile()
	if err != nil {
		return user.User{}, err
	}

	return u, nil
}
