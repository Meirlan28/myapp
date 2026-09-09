package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	users := s.ur.GetAll()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)
}
