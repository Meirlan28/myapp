package main

import (
	"myapp/internal/api"
	"myapp/internal/repository"
)

const port = ":8080"
const fileName = "data.json"

func main() {
	ur := repository.New(fileName)

	server := api.New(ur, port)
	server.Start()
}
