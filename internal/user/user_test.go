package user

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	testName := "test_user"
	testAge := 10
	u, err := NewUser(testName, testAge)
	if err != nil {
		panic(err)
	}
	if u.Name != testName {
		t.Fatalf("Error: name mismatch: %s != %s", u.Name, testName)
		panic(err)
	}
	if u.Age != testAge {
		t.Fatalf("Error: age mismatch: %d != %d", u.Age, testAge)
		panic(err)
	}
	if u.Id == uuid.Nil {
		t.Fatalf("Error: uuid is null: %v", u.Id)
		panic(err)
	}
}

func TestAgeCategory(t *testing.T) {
	tests := []struct {
		user User
		want string
		isError bool
	}{
		{User{uuid.Nil, "junior_1", 0}, Junior, false},
		{User{uuid.Nil, "junior_2", 10}, Junior, false},
		{User{uuid.Nil, "junior_3", 17}, Junior, false},
		{User{uuid.Nil, "middle_1", 18}, Middle, false},
		{User{uuid.Nil, "middle_2", 30}, Middle, false},
		{User{uuid.Nil, "middle_3", 63}, Middle, false},
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
				return
			}
			if err != nil {
				t.Fatalf("failed to calculate age category: %s", err)
			}
			if ag != tt.want {
				t.Fatalf("age category must be %s actual value: %s", tt.want, ag)
			}
		})
	}
}
