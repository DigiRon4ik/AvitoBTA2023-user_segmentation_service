package server

import "user_segmentation_service/internal/server/handlers"

func (api *APIServer) configureRouter() {
	userHandlers := handlers.NewUserHandlers(api.us)
	api.router.HandleFunc("POST /users", userHandlers.CreateHandler)
	api.router.HandleFunc("DELETE /users/{id}", userHandlers.DeleteHandler)
	api.router.HandleFunc("PUT /users/{id}", userHandlers.UpdateHandler)
	api.router.HandleFunc("GET /users/{id}", userHandlers.GetHandler)
	api.router.HandleFunc("GET /users", userHandlers.GetAllHandler)

	segmentHandlers := handlers.NewSegmentHandlers(api.ss)
	api.router.HandleFunc("POST /segments", segmentHandlers.CreateHandler)
	api.router.HandleFunc("DELETE /segments/{slug}", segmentHandlers.DeleteHandler)
	api.router.HandleFunc("PUT /segments/{slug}", segmentHandlers.UpdateHandler)
	api.router.HandleFunc("GET /segments/{slug}", segmentHandlers.GetHandler)
	api.router.HandleFunc("GET /segments", segmentHandlers.GetAllHandler)
}
