// Package server contains the logic for setting up and running an HTTP server to manage users and segments.
// Includes route handling, middleware setup, and server configuration.
package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "user_segmentation_service/api"

	"user_segmentation_service/internal/server/handlers"
)

// configureRouter sets up the HTTP route handlers for users and segments.
func (api *APIServer) configureRouter() {
	api.router.Handle("/swagger/", httpSwagger.WrapHandler)

	userHandler := handlers.NewUserHandler(api.ctx, api.us)
	api.router.HandleFunc("POST /users", userHandler.CreateHandle)
	api.router.HandleFunc("DELETE /users/{id}", userHandler.DeleteHandle)
	api.router.HandleFunc("PUT /users/{id}", userHandler.UpdateHandle)
	api.router.HandleFunc("GET /users/{id}", userHandler.GetHandle)
	api.router.HandleFunc("GET /users", userHandler.GetAllHandle)

	segmentHandler := handlers.NewSegmentHandler(api.ctx, api.ss)
	api.router.HandleFunc("POST /segments", segmentHandler.CreateHandle)
	api.router.HandleFunc("DELETE /segments/{slug}", segmentHandler.DeleteHandle)
	api.router.HandleFunc("PUT /segments/{slug}", segmentHandler.UpdateHandle)
	api.router.HandleFunc("GET /segments/{slug}", segmentHandler.GetHandle)
	api.router.HandleFunc("GET /segments", segmentHandler.GetAllHandle)

	userSegmentsHandler := handlers.NewUserSegmentsHandler(api.ctx, api.uss)
	api.router.HandleFunc("PATCH /users/{id}/segments", userSegmentsHandler.UpdateHandle)
	api.router.HandleFunc("GET /users/{id}/segments", userSegmentsHandler.GetActiveHandle)
	api.router.HandleFunc("GET /users/{id}/segments/history", userSegmentsHandler.GetHistoryCSVHandle)

	fs := http.FileServer(http.Dir("reports"))
	api.router.Handle("/reports/", http.StripPrefix("/reports/", fs))
}
