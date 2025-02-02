// Package server contains the logic for setting up and running an HTTP server to manage users and segments.
// Includes route handling, middleware setup, and server configuration.
package server

import "user_segmentation_service/internal/server/handlers"

// configureRouter sets up the HTTP route handlers for users and segments.
func (api *APIServer) configureRouter() {
	userHandlers := handlers.NewUserHandlers(api.ctx, api.us)
	api.router.HandleFunc("POST /users", userHandlers.CreateHandle)
	api.router.HandleFunc("DELETE /users/{id}", userHandlers.DeleteHandle)
	api.router.HandleFunc("PUT /users/{id}", userHandlers.UpdateHandle)
	api.router.HandleFunc("GET /users/{id}", userHandlers.GetHandle)
	api.router.HandleFunc("GET /users", userHandlers.GetAllHandle)

	segmentHandlers := handlers.NewSegmentHandlers(api.ctx, api.ss)
	api.router.HandleFunc("POST /segments", segmentHandlers.CreateHandle)
	api.router.HandleFunc("DELETE /segments/{slug}", segmentHandlers.DeleteHandle)
	api.router.HandleFunc("PUT /segments/{slug}", segmentHandlers.UpdateHandle)
	api.router.HandleFunc("GET /segments/{slug}", segmentHandlers.GetHandle)
	api.router.HandleFunc("GET /segments", segmentHandlers.GetAllHandle)
}
