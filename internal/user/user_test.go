package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		isError bool
	}{
		{"user_1", 10, false},
		{"user_2", -1, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := NewUser(tt.name, tt.age)
			if tt.isError {
				if err == nil {
					t.Fatalf("must to throw an error u=%v err=%v", u, err)
				}
			} else {
				if err != nil {
					t.Fatalf("failed to create user: %v", err)
				}

				if u.Name != tt.name || u.Age != tt.age || u.Id == uuid.Nil {
					t.Fatalf("Something wrong %v", u)
				}
			}
		})
	}
}

func TestAgeCategory(t *testing.T) {
	tests := []struct {
		user    User
		want    string
		isError bool
	}{
		{User{uuid.Nil, "junior_1", 0}, Underage, false},
		{User{uuid.Nil, "junior_2", 10}, Underage, false},
		{User{uuid.Nil, "junior_3", 17}, Underage, false},
		{User{uuid.Nil, "middle_1", 18}, Adult, false},
		{User{uuid.Nil, "middle_2", 30}, Adult, false},
		{User{uuid.Nil, "middle_3", 63}, Adult, false},
		{User{uuid.Nil, "senior_1", 64}, Senior, false},
		{User{uuid.Nil, "senior_2", 100}, Senior, false},
		{User{uuid.Nil, "senior_3", 120}, Senior, false},
		{User{uuid.Nil, "negative_age", -1}, "", true},
		{User{uuid.Nil, "more_than_max_age", 121}, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.user.Name, func(t *testing.T) {
			ag, err := tt.user.AgeCategory()

			if tt.isError {
				if err == nil {
					t.Fatalf("must to throw an error: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("failed to calculate age category: %s", err)
				}
				if ag != tt.want {
					t.Fatalf("age category must be %s actual value: %s", tt.want, ag)
				}
			}
		})
	}
}
