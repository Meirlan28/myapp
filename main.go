package main

import (
	"myapp/internal/console"
	"myapp/internal/repository"
)

func main() {
	userRepository := repository.New("data.json")
	err := userRepository.Load()
	if err != nil {
		panic(err)
	}

	app := console.New(userRepository)
	app.Start()
}
