package main

import (
	"log/slog"
	"myapp/internal/api"
	"myapp/internal/repository"
	"os"
)

const port = ":8081"
const fileName = "data.json"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	ur := repository.New(fileName, logger)
	err := ur.Load()
	if err != nil {
		panic(err)
	}


	server := api.New(ur, port, logger)
	server.Start()

	// console := console.New(ur)
	// console.Start()
}
