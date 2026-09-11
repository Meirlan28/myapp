package api

import (
	"encoding/json"
	"errors"
	"myapp/internal/repository"
	"myapp/internal/user"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

const DEFAULT_LIMIT = 20
const DEFAULT_OFFSET = 0

func (s *Server) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	min_age_query := r.URL.Query().Get("min_age")
	min_age, err := strconv.Atoi(min_age_query)
	if err != nil {
		if min_age_query == "" {
			min_age = user.MinAge
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
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
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
	}

	limit_query := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limit_query)
	if err != nil {
		if limit_query == "" {
			limit = DEFAULT_LIMIT
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
	}
	offset_query := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offset_query)
	if err != nil {
		if offset_query == "" {
			offset = DEFAULT_OFFSET
		} else {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
	}

	users, err := s.ur.GetAllFiltered(min_age, max_age, limit, offset)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}
	count := s.ur.CountByAge(min_age, max_age)

	var userPage UserPage

	if users == nil {
		userPage = UserPage{[]user.User{}, count, limit, offset}
	} else {
		userPage = UserPage{users, count, limit, offset}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(userPage)
}

func (s *Server) findUserHandler(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	user, err := s.ur.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func (s *Server) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userCreateRequest UserCreateRequest
	err := json.NewDecoder(r.Body).Decode(&userCreateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	if userCreateRequest.Name == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	if len(*userCreateRequest.Name) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	if userCreateRequest.Age == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	if *userCreateRequest.Age < user.MinAge || user.MaxAge < *userCreateRequest.Age {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	u, err := s.ur.Create(*userCreateRequest.Name, *userCreateRequest.Age)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func (s *Server) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userUpdate user.User

	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	var userUpdateRequest UserUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&userUpdateRequest)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	if userUpdateRequest.Name != nil && len(*userUpdateRequest.Name) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	if userUpdateRequest.Age != nil && (*userUpdateRequest.Age < user.MinAge || user.MaxAge < *userUpdateRequest.Age) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ValidationError)
		return
	}

	if userUpdateRequest.Name != nil {
		u, err := s.ur.UpdateName(id, *userUpdateRequest.Name)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
		userUpdate = u
	}

	if userUpdateRequest.Age != nil {
		u, err := s.ur.UpdateAge(id, *userUpdateRequest.Age)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
		userUpdate = u
	}

	userUpdate, err = s.ur.FindByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userUpdate)
}

func (s *Server) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}

	err = s.ur.DeleteByID(id)
	if err != nil {
		if errors.Is(repository.UserNotFound, err) {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{err.Error()})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type UserCreateRequest struct {
	Name *string
	Age  *int
}

type UserUpdateRequest struct {
	Name *string
	Age  *int
}

type UserPage struct {
	Users  []user.User `json:"users"`
	Total  int         `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
