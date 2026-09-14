package api

import (
	"log/slog"
	"myapp/internal/user/transport"
	"net/http"
)

type Server struct {
	mux         *http.ServeMux
	port        string
	Logger      *slog.Logger
	userHandler *transport.UserHandler
}

func New(port string, logger *slog.Logger, userHandler *transport.UserHandler) *Server {
	return &Server{
		http.NewServeMux(),
		port,
		logger,
		userHandler,
	}
}

func (s *Server) Start() {
	s.mux.HandleFunc("GET /users", s.userHandler.GetUsersHandler)
	s.mux.HandleFunc("GET /users/{id}", s.userHandler.FindUserHandler)
	s.mux.HandleFunc("POST /users", s.userHandler.CreateUserHandler)
	s.mux.HandleFunc("PUT /users/{id}", s.userHandler.UpdateUserHandler)
	s.mux.HandleFunc("DELETE /users/{id}", s.userHandler.DeleteUserHandler)
	s.Logger.Info(
		"server started",
		"port", s.port,
	)
	err := http.ListenAndServe(s.port, s.mux)
	if err != nil {
		panic(err)
	}

}
