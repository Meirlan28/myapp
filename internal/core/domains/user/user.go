package user

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
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (u *User) AgeCategory() (string, error) {
	for c := range categories {
		if categories[c].Contains(u.Age) {
			return c, nil
		}
	}
	return "", fmt.Errorf("failed to get age category: Invalid age")
}

func New(name string, age int) (*User, error) {
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
