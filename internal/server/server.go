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

type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) (*models.Segment, error)
	GetByID(ctx context.Context, slug string) (*models.Segment, error)
	GetAll(ctx context.Context) ([]*models.Segment, error)
}

type APIServer struct {
	router *http.ServeMux
	cfg    *Config
	ctx    context.Context
	us     userService
	ss     segmentService
}

func New(ctx context.Context, cfg *Config, us userService, ss segmentService) *APIServer {
	router := http.NewServeMux()

	return &APIServer{
		router: router,
		cfg:    cfg,
		ctx:    ctx,
		us:     us,
		ss:     ss,
	}
}
