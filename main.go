package main

import (
	"myapp/internal/api"
	"myapp/internal/repository"
)

const port = ":8081"
const fileName = "data.json"

func main() {
	ur := repository.New(fileName)
	err := ur.Load()
	if err != nil {
		panic(err)
	}

	server := api.New(ur, port)
	server.Start()

	// console := console.New(ur)
	// console.Start()
}
