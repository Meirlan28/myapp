package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	coreDB "github.com/Meirlan28/myapp/internal/core/database"
	coreLogger "github.com/Meirlan28/myapp/internal/core/logger"
	coreHTTPServer "github.com/Meirlan28/myapp/internal/core/transport/http/server"
	coreHttpServer "github.com/Meirlan28/myapp/internal/core/transport/http/server"
	userRepository "github.com/Meirlan28/myapp/internal/features/user/repository"
	userService "github.com/Meirlan28/myapp/internal/features/user/service"
	userTransportHTTP "github.com/Meirlan28/myapp/internal/features/user/transport/http"
	"github.com/joho/godotenv"
)

func main() {

	// loading env
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println(".env file not found")
	}

	// creating logger
	loggerConfig, err := coreLogger.LoadConfig()
	if err != nil {
		panic(err)
	}
	logger := coreLogger.New(loggerConfig)

	// creating context with gracefull shutdown
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// loading postgres config
	postgresConfig, err := coreDB.LoadConfig()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}

	// connecting to postgres
	db, err := coreDB.NewPostgres(ctx, postgresConfig)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
	defer db.Close()
	logger.Info("connected to postgres")

	serverConfig, err := coreHTTPServer.LoadConfig()
	httpServer := coreHTTPServer.New(serverConfig, logger)

	ur := userRepository.New(db, logger)
	us := userService.New(ur, logger)
	userTransportHTTP := userTransportHTTP.NewUserHTTPHandler(us, logger)

	apiVersionRouterV1 := coreHttpServer.NewAPIVersionRouter(coreHttpServer.ApiVersion1)
	apiVersionRouterV1.AddRoutes(userTransportHTTP.Routes()...)

	httpServer.RegisterAPIRouters(apiVersionRouterV1)
	go httpServer.Start()

	<-ctx.Done()
	logger.Info("shutting down httpServer")

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	err = httpServer.Stop(shutdownCtx)
	if err != nil {
		panic(err)
	}
}
