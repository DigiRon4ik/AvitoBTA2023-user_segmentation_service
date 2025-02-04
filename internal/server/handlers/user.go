// Package handlers provide HTTP request handlers for user segments.
package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"user_segmentation_service/internal/models"
)

// userService defines the methods required for managing users.
type userService interface {
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

// UserHandlers is a structure that contains the user service and context for handling user-related HTTP requests.
type UserHandlers struct {
	users userService
	ctx   context.Context
}

var userHandler = "user handler"

// NewUserHandler creates a new instance of UserHandlers with the provided context and user service.
func NewUserHandler(ctx context.Context, us userService) *UserHandlers {
	return &UserHandlers{
		users: us,
		ctx:   ctx,
	}
}

// CreateHandle handles HTTP POST requests for creating a new user.
//
//	@Summary        Add a user
//	@Description    Creates a user in the database and returns an instance of the user
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Param          User    body        dto.UserCreateRequest    true    "Information about the added user"
//	@Success        201     {object}    dto.UserResponse                 "The user was successfully created"
//	@Router         /users [post]
func (uh *UserHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "CreateHandle"

	var (
		err  error
		user *models.User
	)

	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.users.Create(uh.ctx, user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userHandler, "success", user)
}

// DeleteHandle handles HTTP DELETE requests for deleting a user by ID.
//
//	@Summary        Delete user
//	@Description    Deletes a user from the database
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Param          id      path        int     true    "User ID"
//	@Success        204                                 "The user with this id was successfully deleted"
//	@Router         /users/{id} [delete]
func (uh *UserHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "DeleteHandle"

	var (
		err    error
		userID int
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = uh.users.Delete(uh.ctx, userID); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	slog.Info(fn, "handler", userHandler, "success", userID)
}

// UpdateHandle handles HTTP PUT requests for updating a user by ID.
//
//	@Summary        Update user
//	@Description    Updates the user in the database and returns an instance of the user
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Param          id      path        int                     true    "User ID"
//	@Param          User    body        dto.UserUpdateRequest   true    "User change information"
//	@Success        200     {object}    dto.UserResponse                "A user with this id has been changed"
//	@Router         /users/{id} [put]
func (uh *UserHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "UpdateHandle"

	var (
		err    error
		userID int
		user   *models.User
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.ID = userID
	if err = uh.users.Update(uh.ctx, user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userHandler, "success", user)
}

// GetHandle handles HTTP GET requests for retrieving a user by ID.
//
//	@Summary        Get user
//	@Description    Get user by id
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Param          id      path        int                 true    "User ID"
//	@Success        200     {object}    dto.UserResponse            "A user with this id was received"
//	@Router         /users/{id} [get]
func (uh *UserHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	const fn = "GetHandle"

	var (
		err    error
		userID int
		user   *models.User
	)

	if userID, err = strconv.Atoi(r.PathValue("id")); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user, err = uh.users.GetByID(uh.ctx, userID); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(user); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userHandler, "success", user)
}

// GetAllHandle handles HTTP GET requests for retrieving all users.
//
//	@Summary        Get All users
//	@Description    Get all users from the database
//	@Tags           users
//	@Accept         json
//	@Produce        json
//	@Success        200     {array}     dto.UserResponse        "An array of users was obtained"
//	@Router         /users [get]
func (uh *UserHandlers) GetAllHandle(w http.ResponseWriter, _ *http.Request) {
	const fn = "GetAllHandle"

	var (
		err   error
		users []*models.User
	)

	if users, err = uh.users.GetAll(uh.ctx); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err = json.NewEncoder(w).Encode(users); err != nil {
		slog.Error(fn, "handler", userHandler, "err", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	slog.Info(fn, "handler", userHandler, "success", users)
}
