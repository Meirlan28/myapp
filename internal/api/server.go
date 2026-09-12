package api

import (
	"log/slog"
	"myapp/internal/repository"
	"net/http"
)

type Server struct {
	mux    *http.ServeMux
	port   string
	ur     *repository.UserRepository
	logger *slog.Logger
}

func New(ur *repository.UserRepository, port string, logger *slog.Logger) *Server {
	return &Server{
		http.NewServeMux(),
		port,
		ur,
		logger,
	}
}

func (s *Server) Start() {
	s.mux.HandleFunc("GET /users", s.getUsersHandler)
	s.mux.HandleFunc("GET /users/{id}", s.findUserHandler)
	s.mux.HandleFunc("POST /users", s.createUserHandler)
	s.mux.HandleFunc("PUT /users/{id}", s.updateUserHandler)
	s.mux.HandleFunc("DELETE /users/{id}", s.deleteUserHandler)
	s.logger.Info(
		"server started",
		"port", s.port,
	)
	err := http.ListenAndServe(s.port, s.mux)
	if err != nil {
		panic(err)
	}

}
