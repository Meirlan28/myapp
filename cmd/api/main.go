package main

import (
	"log/slog"
	api "myapp/internal/core/http"
	"myapp/internal/user/repository"
	"myapp/internal/user/service"
	"myapp/internal/user/transport"
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
	us := service.New(ur, logger)
	err := ur.Load()
	if err != nil {
		panic(err)
	}
	userHandler := transport.NewUserHandler(us, logger)

	server := api.New(port, logger, userHandler)
	server.Start()
}
