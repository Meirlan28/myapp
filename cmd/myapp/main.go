package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Meirlan28/myapp/internal/core/logger"
	"github.com/Meirlan28/myapp/internal/core/postgres"
	"github.com/Meirlan28/myapp/internal/core/transport/http/middleware"
	"github.com/Meirlan28/myapp/internal/core/transport/http/server"
	"github.com/Meirlan28/myapp/internal/features/user/repository"
	"github.com/Meirlan28/myapp/internal/features/user/service"
	"github.com/Meirlan28/myapp/internal/features/user/transport/userhttp"
	"github.com/joho/godotenv"
)

func main() {

	// loading env
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println(".env file not found")
	}

	// creating logger
	loggerConfig, err := logger.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger := logger.New(loggerConfig)

	// creating context with gracefull shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// loading postgres config
	postgresConfig, err := postgres.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// connecting to postgres
	db, err := postgres.New(ctx, postgresConfig)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer db.Close()
	logger.Info("connected to postgres")

	serverConfig, err := server.LoadConfig()
	httpServer := server.New(serverConfig, logger,
		middleware.CORS(serverConfig.AllowedOrigins),
		middleware.RequestID(),
		middleware.Logger(logger),
		middleware.Trace(),
		middleware.Recovery())

	ur := repository.New(db, logger)
	us := service.New(ur, logger)
	userTransportHTTP := userhttp.NewHTTPHandler(us, logger)

	apiVersionRouterV1 := server.NewAPIVersionRouter(server.ApiVersion1)
	apiVersionRouterV1.AddRoutes(userTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	if err := httpServer.Start(ctx); err != nil {
		logger.Error("HTTP server run error", "error", err)
	}
}
