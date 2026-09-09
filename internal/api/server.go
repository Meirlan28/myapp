package api

import (
	"fmt"
	"myapp/internal/repository"
	"net/http"
)

type Server struct {
	mux  *http.ServeMux
	port string
	ur   *repository.UserRepository
}

func New(ur *repository.UserRepository, port string) *Server {
	return &Server{
		http.NewServeMux(),
		port,
		ur,
	}
}

func (s *Server) Start() {
	s.mux.HandleFunc("GET /users", s.getUsersHandler)
	s.mux.HandleFunc("GET /users/{id}", s.findUserHandler)
	fmt.Printf("🚀  Application running on http://localhost%s", s.port)
	err := http.ListenAndServe(s.port, s.mux)
	if err != nil {
		panic(err)
	}

}
