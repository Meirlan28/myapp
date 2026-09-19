package user

import (
	"fmt"
)

const MinLimit = 1
const MaxLimit = 100

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

func (u *User) SetAge(age int) error {
	if age < categories[Underage].Min || categories[Senior].Max < age {
		return fmt.Errorf("failed to set age: Invalid age")
	}
	u.Age = age
	return nil
}
