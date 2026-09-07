package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/google/uuid"
)

type UserRepository struct {
	fileName string
	users    []User
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

func (ur *UserRepository) readFromFile() error {
	data, err := os.ReadFile(ur.fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			ur.users = []User{}
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

func (ur *UserRepository) getAll() []User {
	return ur.users
}

func (ur *UserRepository) findById(id uuid.UUID) (User, bool) {
	for i := range ur.users {
		if ur.users[i].Id == id {
			return ur.users[i], true
		}
	}
	return User{}, false
}

func (ur *UserRepository) deleteById(id uuid.UUID) (bool, error) {
	var deleted bool
	for i := range ur.users {
		if ur.users[i].Id == id {
			ur.users = append(ur.users[:i], ur.users[i+1:]...)
			deleted = true
			break
		}
	}
	err := ur.saveToFile()
	if err != nil {
		return deleted, err
	}
	return deleted, nil
}

func (ur *UserRepository) update(id uuid.UUID, name string, age int) (bool, error) {
	var updated bool
	for i, u := range ur.users {
		if u.Id == id {
			err := ur.users[i].SetAge(age)
			if err != nil {
				return false, err
			}
			ur.users[i].Name = name
			updated = true
		}
	}

	err := ur.saveToFile()
	if err != nil {
		return updated, err
	}
	return updated, nil
}

func (ur *UserRepository) create(name string, age int) error {
	user, err := NewUser(name, age)
	if err != nil {
		return err
	}

	ur.users = append(ur.users, user)

	err = ur.saveToFile()
	if err != nil {
		return err
	}

	return nil
}
