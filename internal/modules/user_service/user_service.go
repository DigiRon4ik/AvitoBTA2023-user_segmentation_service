package user_service

import (
	"context"

	"user_segmentation_service/internal/models"
)

type DB interface {
	CreateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, userID int) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]*models.User, error)
}

type UserService struct {
	store DB
}

func NewUserService(store DB) *UserService {
	return &UserService{store: store}
}

func (s *UserService) Create(ctx context.Context, user *models.User) error {
	return s.store.CreateUser(ctx, user)
}

func (s *UserService) Delete(ctx context.Context, userID int) error {
	return s.store.DeleteUser(ctx, userID)
}

func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.store.UpdateUser(ctx, user)
}

func (s *UserService) GetByID(ctx context.Context, userID int) (*models.User, error) {
	return s.store.GetUserByID(ctx, userID)
}

func (s *UserService) GetAll(ctx context.Context) ([]*models.User, error) {
	return s.store.GetAllUsers(ctx)
}
