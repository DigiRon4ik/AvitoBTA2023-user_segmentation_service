// Package server contains the logic for setting up and running an HTTP server to manage users and segments.
// Includes route handling, middleware setup, and server configuration.
package server

import (
	"context"
	"net/http"
	"time"

	"user_segmentation_service/internal/models"
	"user_segmentation_service/internal/server/middlewares"
)

// Config holds configuration values for the API server, such as host and port.
type Config struct {
	Host string `envconfig:"HOST" default:"localhost"`
	Port string `envconfig:"PORT" default:"8080"`
}

// userService defines the methods required for managing users.
type userService interface {
	Create(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, userID int) error
	Update(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, userID int) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
}

// segmentService defines the methods required for managing segments.
type segmentService interface {
	Create(ctx context.Context, seg *models.Segment) error
	Delete(ctx context.Context, slug string) error
	Update(ctx context.Context, seg *models.Segment) error
	GetBySlug(ctx context.Context, slug string) (*models.Segment, error)
	GetAll(ctx context.Context) ([]*models.Segment, error)
}

// APIServer represents the API server, including configuration, router, and services.
type APIServer struct {
	router *http.ServeMux  // HTTP router for handling requests.
	cfg    Config          // Configuration for server settings.
	ctx    context.Context // Application context.
	us     userService     // User service for user-related operations.
	ss     segmentService  // Segment service for segment-related operations.
}

// New creates a new instance of APIServer with the provided context, configuration, and services.
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

// Start begins the HTTP server, listening on the configured host and port.
func (api *APIServer) Start() error {
	api.configureRouter() // Configure the HTTP routes
	server := &http.Server{
		Addr:         api.cfg.Host + ":" + api.cfg.Port,
		Handler:      middlewares.NewMiddleware(api.router), // Apply middleware to the router
		ReadTimeout:  time.Second * 30,                      // Request read timeout
		WriteTimeout: time.Second * 10,                      // Response Record Timeout
		IdleTimeout:  time.Second * 60,                      // Keep-alive connections timeout
	}
	return server.ListenAndServe() // Start the HTTP server
}
