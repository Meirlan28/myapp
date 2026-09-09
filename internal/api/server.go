package api

import (
	"myapp/internal/repository"
	"net/http"
)

type Server struct {
	mux *http.ServeMux
	port string
	ur *repository.UserRepository
}

func New(ur *repository.UserRepository, port string) *Server {
	return &Server{
		http.NewServeMux(),
		port,
		ur,
	}
}

func (s *Server) Start() {
	http.HandleFunc("GET /users", s.usersHandler)
	http.ListenAndServe(s.port, s.mux)
}
