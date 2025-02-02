package server

import "user_segmentation_service/internal/server/handlers"

func (s *APIServer) configureRouter() {
	userHandlers := handlers.NewUserHandlers(s.us)
	s.router.HandleFunc("POST /users", userHandlers.CreateHandler)
	s.router.HandleFunc("DELETE /users/{id}", userHandlers.DeleteHandler)
	s.router.HandleFunc("PUT /users/{id}", userHandlers.UpdateHandler)
	s.router.HandleFunc("GET /users/{id}", userHandlers.GetHandler)
	s.router.HandleFunc("GET /users", userHandlers.GetAllHandler)
}
