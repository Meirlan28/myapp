package console

import (
	"fmt"

	"myapp/internal/repository"
	"myapp/internal/user"

	"github.com/google/uuid"
)

type Console struct {
	ur *repository.UserRepository
}

func New(ur *repository.UserRepository) *Console {
	return &Console{ur: ur}
}

func printUser(u user.User) {
	fmt.Printf("id: %s name: %s age: %d\n", u.Id, u.Name, u.Age)
}

const defaultMenu = `
1 : create user
2 : get user by id
3 : get all users
4 : update user
5 : delete user
exit : exit
`

func (c Console) Start() {
	fmt.Print(defaultMenu)

	for {
		fmt.Print("$ ")

		var input string
		fmt.Scan(&input)

		switch input {
		case "1":
			var name string
			fmt.Printf("Enter user's name: ")
			fmt.Scan(&name)

			var age int
			fmt.Printf("Enter users's age: ")
			fmt.Scan(&age)

			err := c.ur.Create(name, age)
			if err != nil {
				fmt.Print(fmt.Errorf("Error: %w", err))
				break
			}
			fmt.Printf("User created sexsexfully\n")
		case "2":
			fmt.Print("Enter id of the user: ")
			var id uuid.UUID
			var input string
			fmt.Scan(&input)
			id, err := uuid.Parse(input)
			if err != nil {
				fmt.Println("Not valid id")
				continue
			}
			u, found := c.ur.FindByID(id)
			if found {
				fmt.Printf("User found:\n")
				printUser(u)
			} else {
				fmt.Printf("User not found\n")
			}
		case "3":
			users := c.ur.GetAll()
			for _, u := range users {
				printUser(u)
			}
		case "4":
			var id uuid.UUID
			var input string
			fmt.Print("Enter id of the user: ")
			fmt.Scan(&input)
			id, err := uuid.Parse(input)
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %w", err))
				break
			}

			var name string
			fmt.Print("Enter new name: ")
			fmt.Scan(&name)

			var age int
			fmt.Print("Enter new age: ")
			fmt.Scan(&age)

			updated, err := c.ur.Update(id, name, age)
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %w", err))
				break
			}

			if updated {
				fmt.Printf("user updated \n")
			} else {
				fmt.Printf("user with id: %s not found\n", id)
			}
		case "5":
			var id uuid.UUID
			var input string
			fmt.Print("Enter id of the user: ")
			fmt.Scan(&input)
			id, err := uuid.Parse(input)
			if err != nil {
				fmt.Print(fmt.Errorf("Error: %w", err))
				break
			}

			deleted, err := c.ur.DeleteByID(id)
			if err != nil {
				fmt.Println(fmt.Errorf("Error: %w", err))
				break
			}
			if deleted {
				fmt.Printf("User deleted sexsexfully\n")
			} else {
				fmt.Printf("User not found\n")
			}
		case "exit":
			return
		default:
			fmt.Print(defaultMenu)
		}
	}
}
