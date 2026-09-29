package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/Meirlan28/myapp/internal/core/optional"
)

const MaxUserNameLength = 255
const MinUserNameLength = 2

type User struct {
	ID      int    `json:"id"`
	Version int    `json:"version"` // Optimistic locking будет чуть позже
	Name    string `json:"name"`
	Age     int    `json:"age"`
}

func NewUser(name string, age int) (*User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrEmptyName
	}
	if utf8.RuneCountInString(name) < MinUserNameLength || MaxUserNameLength < utf8.RuneCountInString(name) {
		return nil, ErrInvalidName
	}
	if age < MinAge || MaxAge < age {
		return nil, ErrInvalidAge
	}
	return &User{
		Name: name,
		Age:  age,
	}, nil
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrEmptyName
	}
	if u.Age < MinAge || MaxAge < u.Age {
		return ErrInvalidAge
	}
	return nil
}

type UserUpdate struct {
	Name optional.Optional[string]
	Age  *int
}

func NewUserUpdate(
	name optional.Optional[string],
	age *int,
) UserUpdate {
	return UserUpdate{
		Name: name,
		Age:  age,
	}
}

func (p *UserUpdate) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return ErrEmptyName
	}
	return nil
}

func (u *User) ApplyUpdate(patch UserUpdate) error {
	if err := patch.Validate(); err != nil {
		return err
	}

	tmpUser := *u

	if patch.Age != nil {
		tmpUser.Age = *patch.Age
	} else {
		return fmt.Errorf("age is required: %w", ErrInvalidAge)
	}

	if patch.Name.Set {
		tmpUser.Name = *patch.Name.Value
	}

	if err := tmpUser.Validate(); err != nil {
		return err
	}

	*u = tmpUser

	return nil
}
