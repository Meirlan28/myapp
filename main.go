package main

func main() {
	userRepository := UserRepository{
		fileName: "data.json",
		users:    []User{},
	}
	err := userRepository.readFromFile()
	if err != nil {
		panic(err)
	}

	console := Console{ur: &userRepository}

	console.Start()
}
