package handlers

import (
	"context"
	"net/http"

	"user_segmentation_service/internal/models"
)

type userService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type UserHandlers struct {
	users userService
	ctx   context.Context
}

func NewUserHandlers(ctx context.Context, us userService) *UserHandlers {
	return &UserHandlers{
		users: us,
		ctx:   ctx,
	}
}

func (uh *UserHandlers) CreateHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) DeleteHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) UpdateHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) GetHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) GetAllHandle(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
