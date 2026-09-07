package user

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	Id   uuid.UUID
	Name string
	Age  int
}

func NewUser(name string, age int) (User, error) {
	id := uuid.New()
	u := User{
		Id:   id,
		Name: name,
	}
	err := u.SetAge(age)
	if err != nil {
		return User{}, err
	}

	return u, nil
}

func (u User) AgeCategory() (string, error) {
	switch {
	case u.Age >= 0 && u.Age < 18:
		return "underage", nil
	case u.Age >= 18 && u.Age < 64:
		return "adult", nil
	case u.Age >= 64 && u.Age <= 120:
		return "senior", nil
	default:
		return "", errors.New("age must be between 0 and 120")
	}
}

func (u *User) SetAge(age int) error {
	oldAge := u.Age
	u.Age = age
	_, err := u.AgeCategory()
	if err != nil {
		u.Age = oldAge
		return err
	}
	return nil
}
