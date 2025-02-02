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

func (uh *UserHandlers) CreateHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) GetHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func (uh *UserHandlers) GetAllHandler(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
