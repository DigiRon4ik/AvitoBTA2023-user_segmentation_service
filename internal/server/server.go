package server

import (
	"context"
	"net/http"
	"time"

	"user_segmentation_service/internal/models"
)

type Config struct {
	Host string `envconfig:"HOST" default:"localhost"`
	Port string `envconfig:"PORT" default:"8080"`
}

type userService interface {
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) error
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) error
	GetBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAll(ctx context.Context) ([]*models.Segment, error)
}

type APIServer struct {
	router *http.ServeMux
	cfg    Config
	ctx    context.Context
	us     userService
	ss     segmentService
}

func New(ctx context.Context, cfg Config, us userService, ss segmentService) *APIServer {
	router := http.NewServeMux()

	return &APIServer{
		router: router,
		cfg:    cfg,
		ctx:    ctx,
		us:     us,
		ss:     ss,
	}
}

func (api *APIServer) Start() error {
	api.configureRouter()
	server := &http.Server{
		Addr:         api.cfg.Host + ":" + api.cfg.Port,
		Handler:      api.router,
		ReadTimeout:  time.Second * 30, // Request read timeout
		WriteTimeout: time.Second * 10, // Response Record Timeout
		IdleTimeout:  time.Second * 60, // Keep-alive connections timeout
	}
	return server.ListenAndServe()
}
