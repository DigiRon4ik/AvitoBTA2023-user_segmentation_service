// Package dto for Swagger
package dto

import "time"

// UserCreateRequest for Swagger
//
//	@Description User information on creation
type UserCreateRequest struct {
	// required: true
	Name string `json:"name"`
}

// UserUpdateRequest for Swagger
//
//	@Description User information when updating
type UserUpdateRequest struct {
	// required: true
	Name string `json:"name"`
}

// UserResponse for Swagger
//
//	@Description User information at creation/update
type UserResponse struct {
	// read only: true
	ID   int    `json:"id"`
	Name string `json:"name"`
	// read only: true
	CreatedAt time.Time `json:"created_at"`
}
