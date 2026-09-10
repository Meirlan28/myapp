package api

import (
	"encoding/json"
	"fmt"
	"myapp/internal/user"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	min_age_query := r.URL.Query().Get("min_age")
	min_age, err := strconv.Atoi(min_age_query)
	if err != nil {
		if min_age_query == "" {
			min_age = user.MinAge
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Sprintf("Error: invalid age: %v", err))
			return
		}
	}
	max_age_query := r.URL.Query().Get("max_age")
	max_age, err := strconv.Atoi(max_age_query)
	if err != nil {
		if max_age_query == "" {
			max_age = user.MaxAge
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(fmt.Sprintf("Error: invalid age: %v", err))
			return
		}
	}
	users := s.ur.GetAllFiltered(min_age, max_age)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(users)
}

func (s *Server) findUserHandler(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %s is not valid uuid", pathId))
		return
	}

	user, err := s.ur.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode("Error: user not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userCreateRequest UserCreateRequest
	json.NewDecoder(r.Body).Decode(&userCreateRequest)

	u, err := s.ur.Create(userCreateRequest.Name, userCreateRequest.Age)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %s is not valid uuid", pathId))
		return
	}

	var userUpdateRequest UserUpdateRequest
	json.NewDecoder(r.Body).Decode(&userUpdateRequest)

	u, err := s.ur.Update(id, userUpdateRequest.Name, userUpdateRequest.Age)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %v", err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %s is not valid uuid", pathId))
		return
	}

	err = s.ur.DeleteByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(fmt.Sprintf("Error: %v", err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type UserCreateRequest struct {
	Name string
	Age  int
}

type UserUpdateRequest struct {
	Name string
	Age  int
}
