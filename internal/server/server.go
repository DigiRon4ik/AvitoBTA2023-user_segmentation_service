package server

import (
	"context"
	"net/http"

	"user_segmentation_service/internal/models"
)

type Config struct {
	Host string `envconfig:"HOST" default:"localhost"`
	Port string `envconfig:"PORT" default:"8080"`
}

type userService interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) (*models.User, error)
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type APIServer struct {
	router *http.ServeMux
	cfg    *Config
	ctx    context.Context
	us     userService
}

func New(ctx context.Context, cfg *Config, us userService) *APIServer {
	router := http.NewServeMux()

	return &APIServer{
		router: router,
		cfg:    cfg,
		ctx:    ctx,
		us:     us,
	}
}
