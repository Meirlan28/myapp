package transport

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"myapp/internal/core/apperrors"
	"myapp/internal/core/domains/user"
	"myapp/internal/user/service"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

const DefaultLimit = 20
const DefaultOffset = 0

type UserHandler struct {
	Us     *service.UserService
	Logger *slog.Logger
}

func NewUserHandler(us *service.UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{us, logger}
}

func (uh *UserHandler) GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var appErr apperrors.AppError
	minAgeQuery := r.URL.Query().Get("min_age")
	minAge, err := strconv.Atoi(minAgeQuery)
	if err != nil {
		if minAgeQuery == "" {
			minAge = user.MinAge
		} else {
			uh.SendError(w, apperrors.NewBadRequestError(err))
			return
		}
	}

	maxAgeQuery := r.URL.Query().Get("max_age")
	maxAge, err := strconv.Atoi(maxAgeQuery)
	if err != nil {
		if maxAgeQuery == "" {
			maxAge = user.MaxAge
		} else {
			uh.SendError(w, apperrors.NewBadRequestError(err))
			return
		}
	}

	limitQuery := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitQuery)
	if err != nil {
		if limitQuery == "" {
			limit = DefaultLimit
		} else {
			uh.SendError(w, apperrors.NewBadRequestError(err))
			return
		}
	}
	offsetQuery := r.URL.Query().Get("offset")
	offset, err := strconv.Atoi(offsetQuery)
	if err != nil {
		if offsetQuery == "" {
			offset = DefaultOffset
		} else {
			uh.SendError(w, apperrors.NewBadRequestError(err))
			return
		}
	}

	uh.Logger.Info("FindAll",
		"min_age", minAge,
		"max_age", maxAge,
		"limit", limit,
		"offset", offset)

	users, appErr := uh.Us.FindAll(minAge, maxAge, limit, offset)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}
	count, appErr := uh.Us.CountByAge(minAge, maxAge)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}

	var userPage UserPage

	if users == nil {
		userPage = UserPage{[]user.User{}, count, limit, offset}
	} else {
		userPage = UserPage{users, count, limit, offset}
	}

	uh.SendResponse(w, userPage, http.StatusOK)
}

func (uh *UserHandler) FindUserHandler(w http.ResponseWriter, r *http.Request) {
	var appErr apperrors.AppError
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)

	if err != nil {
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	u, appErr := uh.Us.FindById(id)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}

	uh.SendResponse(w, u, http.StatusOK)
}

func (uh *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var appErr apperrors.AppError
	var userCreateRequest UserCreateRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&userCreateRequest)
	if err != nil {
		uh.Logger.Info("failed to parsed body",
			"body", userCreateRequest, "error", err)
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	var extra any
	err = decoder.Decode(&extra)
	if !errors.Is(err, io.EOF) {
		uh.Logger.Info(
			"request body contains extra data",
			"error", err,
		)

		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	if userCreateRequest.Name == nil {
		uh.SendError(w, apperrors.NewBadRequestError(InvalidNameError))
		return
	}

	*userCreateRequest.Name = strings.TrimSpace(*userCreateRequest.Name)

	if utf8.RuneCountInString(*userCreateRequest.Name) < 2 {
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	if userCreateRequest.Age == nil {
		uh.SendError(w, apperrors.NewBadRequestError(InvalidAgeError))
		return
	}

	if *userCreateRequest.Age < user.MinAge || user.MaxAge < *userCreateRequest.Age {
		uh.SendError(w, apperrors.NewBadRequestError(InvalidAgeError))
		return
	}

	u, appErr := uh.Us.Create(*userCreateRequest.Name, *userCreateRequest.Age)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}

	uh.SendResponse(w, u, http.StatusCreated)
}

func (uh *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	var appErr apperrors.AppError

	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		uh.Logger.Info("failed to parse id",
			"id", id,
			"error", err)
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}
	uh.Logger.Info("id parsed",
		"id", id)

	var userUpdateRequest UserUpdateRequest
	err = json.NewDecoder(r.Body).Decode(&userUpdateRequest)
	if err != nil {
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	u, AppErr := uh.Us.FindById(id)
	if AppErr != nil {
		uh.SendError(w, AppErr)
		return
	}

	if userUpdateRequest.Name == nil {
		userUpdateRequest.Name = &u.Name
	}

	if userUpdateRequest.Age == nil {
		userUpdateRequest.Age = &u.Age
	}

	*userUpdateRequest.Name = strings.TrimSpace(*userUpdateRequest.Name)

	if utf8.RuneCountInString(*userUpdateRequest.Name) < 2 {
		uh.SendError(w, apperrors.NewBadRequestError(InvalidAgeError))
		return
	}

	if *userUpdateRequest.Age < user.MinAge || user.MaxAge < *userUpdateRequest.Age {
		uh.SendError(w, apperrors.NewBadRequestError(InvalidAgeError))
		return
	}

	u, appErr = uh.Us.Update(id, *userUpdateRequest.Name, *userUpdateRequest.Age)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}

	uh.SendResponse(w, u, http.StatusOK)
}

func (uh *UserHandler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var appErr apperrors.AppError
	pathId := r.PathValue("id")
	id, err := uuid.Parse(pathId)
	if err != nil {
		uh.SendError(w, apperrors.NewBadRequestError(err))
		return
	}

	appErr = uh.Us.Delete(id)
	if appErr != nil {
		uh.SendError(w, appErr)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (uh *UserHandler) SendError(w http.ResponseWriter, err apperrors.AppError) {
	encoder := json.NewEncoder(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.GetCode())

	_ = encoder.Encode(ErrorResponse{err.Error()})
}

func (uh *UserHandler) SendResponse(w http.ResponseWriter, response any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}
