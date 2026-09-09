package user

import (
	"fmt"
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
	for c := range categories {
		if categories[c].Contains(u.Age) {
			return c, nil
		}
	}
	return "", fmt.Errorf("Invalid age")
}

func (u *User) SetAge(age int) error {
	if age < categories[Underage].Min || categories[Underage].Max < age {
		return fmt.Errorf("failed to set age: Invalid age")
	}
	u.Age = age
	return nil
}
