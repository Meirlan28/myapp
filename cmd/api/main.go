package main

import (
	"context"
	"log/slog"
	api "myapp/internal/core/http"
	"myapp/internal/user/repository"
	"myapp/internal/user/service"
	"myapp/internal/user/transport"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

const port = ":8081"

func main() {
	logger := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				Level: slog.LevelInfo,
			},
		),
	)

	if err := godotenv.Load(".env.local"); err != nil {
		logger.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	databaseURL := os.Getenv("DATABASE_URL")

	ctx := context.Background()

	db, err := pgxpool.New(
		ctx,
		databaseURL,
	)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	logger.Info("connected to postgres")

	ur := repository.New(db, logger)
	us := service.New(ur, logger)
	userHandler := transport.NewUserHandler(us, logger)

	server := api.New(port, logger, userHandler)
	server.Start()
}
