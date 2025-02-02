// Package user_service provides business logic for managing users.
package user_service

import (
	"context"

	"user_segmentation_service/internal/models"
)

// DB defines the required database operations for user management.
type DB interface {
	CreateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

// UserService handles operations related to users.
type UserService struct {
	store DB
}

// NewUserService creates a new instance of UserService.
func NewUserService(store DB) *UserService {
	return &UserService{store: store}
}

// Create adds a new user to the database.
func (s *UserService) Create(ctx context.Context, user *models.User) error {
	return s.store.CreateUser(ctx, user)
}

// Delete removes a user by ID.
func (s *UserService) Delete(ctx context.Context, userID int) error {
	return s.store.DeleteUser(ctx, userID)
}

// Update modifies an existing user.
func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.store.UpdateUser(ctx, user)
}

// GetByID retrieves a user by ID.
func (s *UserService) GetByID(ctx context.Context, userID int) (*models.User, error) {
	return s.store.GetUserByID(ctx, userID)
}

// GetAll returns all users from the database.
func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	return s.store.GetAllUsers(ctx)
}
