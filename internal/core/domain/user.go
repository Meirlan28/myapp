package domain

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const MinLimit = 1
const MaxLimit = 100

const MaxUserNameLength = 255
const MinUserNameLength = 2

type User struct {
	ID          int    `json:"id"`
	Version     int    `json:"version"`
	Name        string `json:"name"`
	Age         int    `json:"age"`
	PhoneNumber string `json:"phone_number"`
}

func (u *User) AgeCategory() (string, error) {
	for c := range categories {
		if categories[c].Contains(u.Age) {
			return c, nil
		}
	}
	return "", fmt.Errorf("failed to get age category: Invalid age")
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
	Name Nullable[string]
	Age  Nullable[int]
}

func NewUserUpdate(
	name Nullable[string],
	age Nullable[int],
) UserUpdate {
	return UserUpdate{
		Name: name,
		Age:  age,
	}
}

func (p *UserUpdate) Validate() error {
	if p.Name.Set && p.Name.Value == nil {
		return fmt.Errorf(
			"`Name` can't be updated to NULL: %w",
			fmt.Errorf("`Name` can't be updated to NULL"),
		)
	}

	return nil
}

func (u *User) ApplyUpdate(patch UserUpdate) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate user patch: %w", err)
	}

	tmp := *u

	if patch.Name.Set {
		tmp.Name = *patch.Name.Value
	}

	if patch.Age.Set {
		tmp.Age = *patch.Age.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp

	return nil
}
